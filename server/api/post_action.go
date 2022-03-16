// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/mscalendar"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/mscalendar/views"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
)

func (api *api) preprocessAction(w http.ResponseWriter, req *http.Request) (mscal mscalendar.MSCalendar, user *mscalendar.User, eventID string, option string, postID string) {
	mattermostUserID := req.Header.Get("Mattermost-User-ID")
	if mattermostUserID == "" {
		utils.SlackAttachmentError(w, "Error: not authorized")
		return nil, nil, "", "", ""
	}

	request := utils.PostActionIntegrationRequestFromJson(req.Body)
	if request == nil {
		utils.SlackAttachmentError(w, "Error: invalid request")
		return nil, nil, "", "", ""
	}

	eventID, ok := request.Context[config.EventIDKey].(string)
	if !ok {
		utils.SlackAttachmentError(w, "Error: missing event ID")
		return nil, nil, "", "", ""
	}
	option, _ = request.Context["selected_option"].(string)
	mscal = mscalendar.New(api.Env, mattermostUserID)

	return mscal, mscalendar.NewUser(mattermostUserID), eventID, option, request.PostId
}

func (api *api) postActionAccept(w http.ResponseWriter, req *http.Request) {
	mscalendar, user, eventID, _, _ := api.preprocessAction(w, req)
	if eventID == "" {
		return
	}
	err := mscalendar.AcceptEvent(user, eventID)
	if err != nil {
		api.Logger.Warnf("Failed to accept event. err=%v", err)
		utils.SlackAttachmentError(w, "Error: Failed to accept event: "+err.Error())
		return
	}
}

func (api *api) postActionDecline(w http.ResponseWriter, req *http.Request) {
	mscalendar, user, eventID, _, _ := api.preprocessAction(w, req)
	if eventID == "" {
		return
	}
	err := mscalendar.DeclineEvent(user, eventID)
	if err != nil {
		utils.SlackAttachmentError(w, "Error: Failed to decline event: "+err.Error())
		return
	}
}

func (api *api) postActionTentative(w http.ResponseWriter, req *http.Request) {
	mscalendar, user, eventID, _, _ := api.preprocessAction(w, req)
	if eventID == "" {
		return
	}
	err := mscalendar.TentativelyAcceptEvent(user, eventID)
	if err != nil {
		utils.SlackAttachmentError(w, "Error: Failed to tentatively accept event: "+err.Error())
		return
	}
}

func (api *api) postActionRespond(w http.ResponseWriter, req *http.Request) {
	calendar, user, eventID, option, postID := api.preprocessAction(w, req)
	if eventID == "" {
		return
	}
	err := calendar.RespondToEvent(user, eventID, option)
	if err != nil && !isAcceptedError(err) && !isNotFoundError(err) && !isCanceledError(err) {
		utils.SlackAttachmentError(w, "Error: Failed to respond to event: "+err.Error())
		return
	}

	if err != nil && isCanceledError(err) {
		utils.SlackAttachmentError(w, "Error: Cannot respond to the event because it is already canceled.")
		return
	}

	p, appErr := api.PluginAPI.GetPost(postID)
	if appErr != nil {
		utils.SlackAttachmentError(w, "Error: Failed to update the post: "+appErr.Error())
		return
	}

	sas := p.Attachments()
	if len(sas) == 0 {
		utils.SlackAttachmentError(w, "Error: Failed to update the post: No attachments found")
		return
	}

	sa := sas[0]

	if err == nil || isAcceptedError(err) {
		sa.Fields = append(sa.Fields, &model.SlackAttachmentField{
			Title: "Response",
			Value: fmt.Sprintf("You have %s this event", prettyOption(option)),
			Short: false,
		})
	}

	sa.Actions = []*model.PostAction{}
	postResponse := model.PostActionIntegrationResponse{}
	model.ParseSlackAttachment(p, []*model.SlackAttachment{sa})

	postResponse.Update = p

	if err != nil && isNotFoundError(err) {
		postResponse.EphemeralText = "Event has changed since this message. Please change your status directly on MS Calendar."
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(utils.ResponseToJson(postResponse))
}

func prettyOption(option string) string {
	switch option {
	case mscalendar.OptionYes:
		return "accepted"
	case mscalendar.OptionNo:
		return "declined"
	case mscalendar.OptionMaybe:
		return "tentatively accepted"
	default:
		return ""
	}
}

func (api *api) postActionConfirmStatusChange(w http.ResponseWriter, req *http.Request) {
	mattermostUserID := req.Header.Get("Mattermost-User-ID")
	if mattermostUserID == "" {
		utils.SlackAttachmentError(w, "Not authorized.")
		return
	}

	response := model.PostActionIntegrationResponse{}
	post := &model.Post{}

	request := utils.PostActionIntegrationRequestFromJson(req.Body)
	if request == nil {
		utils.SlackAttachmentError(w, "Invalid request.")
		return
	}

	value, ok := request.Context["value"].(bool)
	if !ok {
		utils.SlackAttachmentError(w, `No recognizable value for property "value".`)
		return
	}

	returnText := "The status has not been changed."
	if value {
		changeTo, ok := request.Context["change_to"]
		if !ok {
			utils.SlackAttachmentError(w, "No state to change to.")
			return
		}
		stringChangeTo := changeTo.(string)
		prettyChangeTo, ok := request.Context["pretty_change_to"]
		if !ok {
			prettyChangeTo = changeTo
		}
		stringPrettyChangeTo := prettyChangeTo.(string)

		status, err := api.PluginAPI.GetMattermostUserStatus(mattermostUserID)
		if err != nil {
			utils.SlackAttachmentError(w, "Cannot get current status.")
			api.Logger.Debugf("cannot get user status, err=%s", err)
			return
		}

		user, err := api.Store.LoadUser(mattermostUserID)
		if err != nil {
			utils.SlackAttachmentError(w, "Cannot load user")
			return
		}

		user.LastStatus = ""
		if status.Manual {
			user.LastStatus = status.Status
		}

		// Handle custom status change for user
		api.handleCustomStatusChange(w, request, user)

		err = api.Store.StoreUser(user)
		if err != nil {
			utils.SlackAttachmentError(w, "Cannot update user")
		}
		_, _ = api.PluginAPI.UpdateMattermostUserStatus(mattermostUserID, stringChangeTo)
		returnText = fmt.Sprintf("The status has been changed to %s.", stringPrettyChangeTo)
	}

	eventInfo, err := getEventInfo(request.Context)
	if err != nil {
		utils.SlackAttachmentError(w, err.Error())
		return
	}

	if eventInfo != "" {
		returnText = eventInfo + "\n" + returnText
	}

	sa := &model.SlackAttachment{
		Title:    "Status Change",
		Text:     returnText,
		Fallback: "Status Change: " + returnText,
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{sa})

	response.Update = post
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(utils.ResponseToJson(response))
}

func getEventInfo(ctx map[string]interface{}) (string, error) {
	hasEvent, ok := ctx["hasEvent"].(bool)
	if !ok {
		return "", errors.New("cannot check whether there is an event attached")
	}
	if !hasEvent {
		return "", nil
	}

	subject, ok := ctx["subject"].(string)
	if !ok {
		return "", errors.New("cannot find the event subject")
	}

	weblink, ok := ctx["weblink"].(string)
	if !ok {
		return "", errors.New("cannot find the event weblink")
	}

	marshalledStartTime, ok := ctx["startTime"].(string)
	if !ok {
		return "", errors.New("cannot find the event start time")
	}
	var startTime time.Time
	err := json.Unmarshal([]byte(marshalledStartTime), &startTime)
	if err != nil {
		return "", fmt.Errorf("error occurred while unmarshalling start time. Error: %s", err.Error())
	}

	return views.RenderEventWillStartLine(subject, weblink, startTime), nil
}

func isAcceptedError(err error) bool {
	return strings.Contains(err.Error(), "202 Accepted")
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "404 Not Found")
}

func isCanceledError(err error) bool {
	return strings.Contains(err.Error(), "You can't respond to a meeting that's been canceled.")
}

func (api *api) handleCustomStatusChange(w http.ResponseWriter, request *model.PostActionIntegrationRequest, user *store.User) {
	currentCustomStatus, err := api.PluginAPI.GetMattermostUserCustomStatus(user.MattermostUserID)
	if err != nil {
		utils.SlackAttachmentError(w, "Cannot get user custom status.")
		api.Logger.Debugf("cannot get user custom status, err=%s", err)
		return
	}

	removeCustomStatus, ok := request.Context["removeCustomStatus"]
	if !ok {
		utils.SlackAttachmentError(w, `No recognizable value for property "removeCustomStatus".`)
		return
	}
	if removeCustomStatus.(bool) {
		if currentCustomStatus != nil && currentCustomStatus.Text == config.CustomStatusText {
			_ = api.PluginAPI.RemoveMattermostUserCustomStatus(user.MattermostUserID)
		}
		return
	}
	customStatusText, ok := request.Context["customStatusText"]
	if !ok {
		utils.SlackAttachmentError(w, `No recognizable value for property "customStatusText".`)
		return
	}
	customStatusEmoji, ok := request.Context["customStatusEmoji"]
	if !ok {
		utils.SlackAttachmentError(w, `No recognizable value for property "customStatusEmoji".`)
		return
	}
	customStatusExpiresAt, ok := request.Context["customStatusExpiresAt"]
	if !ok {
		utils.SlackAttachmentError(w, `No recognizable value for property "customStatusExpiresAt".`)
		return
	}
	customStatusParsedExpiresAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", customStatusExpiresAt.(string))
	if err != nil {
		utils.SlackAttachmentError(w, `error while parsing custom status expiresAt value".`)
		return
	}
	customStatusDuration, ok := request.Context["customStatusDuration"]
	if !ok {
		utils.SlackAttachmentError(w, `No recognizable value for property "customStatusDuration".`)
		return
	}

	user.LastCustomStatus = nil
	if currentCustomStatus != nil && currentCustomStatus.Text != config.CustomStatusText {
		user.LastCustomStatus = currentCustomStatus
	}
	customStatus := &model.CustomStatus{
		Text:      customStatusText.(string),
		Emoji:     customStatusEmoji.(string),
		Duration:  customStatusDuration.(string),
		ExpiresAt: customStatusParsedExpiresAt,
	}
	_ = api.PluginAPI.UpdateMattermostUserCustomStatus(user.MattermostUserID, customStatus)
}
