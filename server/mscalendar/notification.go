// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
//  See License for license information.

package mscalendar

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/mscalendar/views"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/fields"
)

const maxQueueSize = 1024

const (
	FieldSubject        = "Subject"
	FieldBodyPreview    = "BodyPreview"
	FieldImportance     = "Importance"
	FieldDuration       = "Duration"
	FieldWhen           = "When"
	FieldLocation       = "Location"
	FieldAttendees      = "Attendees"
	FieldOrganizer      = "Organizer"
	FieldResponseStatus = "ResponseStatus"
)

const (
	OptionYes          = "Yes"
	OptionNotResponded = "Not responded"
	OptionNo           = "No"
	OptionMaybe        = "Maybe"
)

const (
	ResponseYes   = "accepted"
	ResponseMaybe = "tentativelyAccepted"
	ResponseNo    = "declined"
	ResponseNone  = "notResponded"
)

var notificationFieldOrder []string = []string{
	FieldWhen,
	FieldLocation,
	FieldAttendees,
	FieldImportance,
}

type NotificationProcessor interface {
	Configure(Env)
	Enqueue(notifications ...*remote.Notification) error
	Quit()
}

type notificationProcessor struct {
	Env
	envChan chan Env

	queue chan *remote.Notification
	quit  chan bool
}

func NewNotificationProcessor(env Env) NotificationProcessor {
	processor := &notificationProcessor{
		Env:     env,
		envChan: make(chan (Env)),
		queue:   make(chan (*remote.Notification), maxQueueSize),
		quit:    make(chan (bool)),
	}
	go processor.work()
	return processor
}

func (processor *notificationProcessor) Enqueue(notifications ...*remote.Notification) error {
	for _, n := range notifications {
		select {
		case processor.queue <- n:
		default:
			return fmt.Errorf("webhook notification: queue full, dropped notification")
		}
	}
	return nil
}

func (processor *notificationProcessor) Configure(env Env) {
	processor.envChan <- env
}

func (processor *notificationProcessor) Quit() {
	processor.quit <- true
}

func (processor *notificationProcessor) work() {
	for {
		select {
		case n := <-processor.queue:
			err := processor.processNotification(n)
			if err != nil && err != store.ErrNotFound {
				processor.Logger.With(bot.LogContext{
					"subscriptionID": n.SubscriptionID,
				}).Infof("webhook notification: failed: `%v`.", err)
			}

		case env := <-processor.envChan:
			processor.Env = env

		case <-processor.quit:
			return
		}
	}
}

func (processor *notificationProcessor) processNotification(n *remote.Notification) error {
	sub, err := processor.Store.LoadSubscription(n.SubscriptionID)
	if err != nil {
		return err
	}
	creator, err := processor.Store.LoadUser(sub.MattermostCreatorID)
	if err != nil {
		return err
	}
	client := processor.Remote.MakeClient(context.Background())
	eventData, err := client.GetNotificationData(creator.Remote.Mail, n.EventID, n.SubscriptionID)
	if err != nil {
		return err
	}
	user, userErr := processor.Env.PluginAPI.GetMattermostUser(sub.MattermostCreatorID)
	if userErr != nil {
		return userErr
	}
	// Setting user local timezone
	eventData.Event.TimeZone = model.GetPreferredTimezone(user.Timezone)
	sa := processor.newEventSlackAttachment(eventData)
	_, err = processor.Poster.DMWithAttachments(creator.MattermostUserID, sa)
	if err != nil {
		return err
	}

	return nil
}

func (processor *notificationProcessor) newSlackAttachment(n *remote.Notification) *model.SlackAttachment {
	title := views.EnsureSubject(n.Event.Subject)
	titleLink := n.Event.Weblink
	return &model.SlackAttachment{
		AuthorName: n.Event.Organizer.EmailAddress.Name,
		AuthorLink: "mailto:" + n.Event.Organizer.EmailAddress.Address,
		TitleLink:  titleLink,
		Title:      title,
		Fallback:   fmt.Sprintf("[%s](%s)", title, titleLink),
	}
}

func (processor *notificationProcessor) newEventSlackAttachment(n *remote.Notification) *model.SlackAttachment {
	sa := processor.newSlackAttachment(n)
	fields := eventToFields(n.Event)
	for _, k := range notificationFieldOrder {
		v := fields[k]

		sa.Fields = append(sa.Fields, &model.SlackAttachmentField{
			Title: k,
			Value: strings.Join(v.Strings(), ", "),
			Short: true,
		})
	}

	if n.Event.ResponseRequested && !n.Event.IsOrganizer {
		sa.Actions = NewPostActionForEventResponse(n.Event.ID, n.Event.ResponseStatus.Response, processor.actionURL(config.PathRespond))
	}
	return sa
}

// TODO: remove if not required after completing the entire flow of subscription and notifications
// func (processor *notificationProcessor) updatedEventSlackAttachment(n *remote.Notification, prior *remote.Event, timezone string) (bool, *model.SlackAttachment) {
// 	sa := processor.newSlackAttachment(n)
// 	sa.Title = "(updated) " + sa.Title

// 	newFields := eventToFields(n.Event, timezone)
// 	priorFields := eventToFields(prior, timezone)
// 	changed, added, updated, deleted := fields.Diff(priorFields, newFields)
// 	if !changed {
// 		return false, nil
// 	}

// 	allChanges := append(added, updated...)
// 	allChanges = append(allChanges, deleted...)

// 	hasImportantChanges := false
// 	for _, k := range allChanges {
// 		if isImportantChange(k) {
// 			hasImportantChanges = true
// 			break
// 		}
// 	}

// 	if !hasImportantChanges {
// 		return false, nil
// 	}

// 	for _, k := range added {
// 		if !isImportantChange(k) {
// 			continue
// 		}
// 		sa.Fields = append(sa.Fields, &model.SlackAttachmentField{
// 			Title: k,
// 			Value: strings.Join(newFields[k].Strings(), ", "),
// 			Short: true,
// 		})
// 	}
// 	for _, k := range updated {
// 		if !isImportantChange(k) {
// 			continue
// 		}
// 		sa.Fields = append(sa.Fields, &model.SlackAttachmentField{
// 			Title: k,
// 			Value: fmt.Sprintf("~~%s~~ \u2192 %s", strings.Join(priorFields[k].Strings(), ", "), strings.Join(newFields[k].Strings(), ", ")),
// 			Short: true,
// 		})
// 	}
// 	for _, k := range deleted {
// 		if !isImportantChange(k) {
// 			continue
// 		}
// 		sa.Fields = append(sa.Fields, &model.SlackAttachmentField{
// 			Title: k,
// 			Value: fmt.Sprintf("~~%s~~", strings.Join(priorFields[k].Strings(), ", ")),
// 			Short: true,
// 		})
// 	}

// 	if n.Event.ResponseRequested && !n.Event.IsOrganizer && !n.Event.IsCancelled {
// 		sa.Actions = NewPostActionForEventResponse(n.Event.ID, n.Event.ResponseStatus.Response, processor.actionURL(config.PathRespond))
// 	}
// 	return true, sa
// }

// func isImportantChange(fieldName string) bool {
// 	for _, ic := range importantNotificationChanges {
// 		if ic == fieldName {
// 			return true
// 		}
// 	}
// 	return false
// }

func (processor *notificationProcessor) actionURL(action string) string {
	return fmt.Sprintf("%s%s%s", processor.Config.PluginURLPath, config.PathPostAction, action)
}

func NewPostActionForEventResponse(eventID, response, url string) []*model.PostAction {
	context := map[string]interface{}{
		config.EventIDKey: eventID,
	}

	pa := &model.PostAction{
		Name: "Response",
		Type: model.POST_ACTION_TYPE_SELECT,
		Integration: &model.PostActionIntegration{
			URL:     url,
			Context: context,
		},
	}

	for _, o := range []string{OptionNotResponded, OptionYes, OptionNo, OptionMaybe} {
		pa.Options = append(pa.Options, &model.PostActionOptions{Text: o, Value: o})
	}
	switch response {
	case ResponseNone:
		pa.DefaultOption = OptionNotResponded
	case ResponseYes:
		pa.DefaultOption = OptionYes
	case ResponseNo:
		pa.DefaultOption = OptionNo
	case ResponseMaybe:
		pa.DefaultOption = OptionMaybe
	}
	return []*model.PostAction{pa}
}

func eventToFields(e *remote.Event) fields.Fields {
	date := func(dtStart, dtEnd *remote.DateTime) (time.Time, time.Time, string) {
		if dtStart == nil || dtEnd == nil {
			return time.Time{}, time.Time{}, "n/a"
		}

		dtStart = dtStart.In(e.TimeZone)
		dtEnd = dtEnd.In(e.TimeZone)
		tStart := dtStart.Time()
		tEnd := dtEnd.Time()
		startFormat := "Monday, January 02"
		if tStart.Year() != time.Now().Year() {
			startFormat = "Monday, January 02, 2006"
		}
		startFormat += " · (" + time.Kitchen
		endFormat := " - " + time.Kitchen + ")"
		return tStart, tEnd, tStart.Format(startFormat) + tEnd.Format(endFormat)
	}

	start, end, formattedDate := date(e.Start, e.End)

	minutes := int(end.Sub(start).Round(time.Minute).Minutes())
	hours := int(end.Sub(start).Hours())
	minutes -= hours * 60
	days := int(end.Sub(start).Hours()) / 24
	hours -= days * 24

	dur := ""
	switch {
	case days > 0:
		dur = fmt.Sprintf("%v days", days)

	case e.IsAllDay:
		dur = "all-day"

	default:
		switch hours {
		case 0:
			// ignore
		case 1:
			dur = "one hour"
		default:
			dur = fmt.Sprintf("%v hours", hours)
		}
		if minutes > 0 {
			if dur != "" {
				dur += ", "
			}
			dur += fmt.Sprintf("%v minutes", minutes)
		}
	}

	attendees := []fields.Value{}
	for _, a := range e.Attendees {
		attendees = append(attendees, fields.NewStringValue(
			fmt.Sprintf("[%s](mailto:%s)",
				a.EmailAddress.Name, a.EmailAddress.Address)))
	}

	if len(attendees) == 0 {
		attendees = append(attendees, fields.NewStringValue("None"))
	}

	ff := fields.Fields{
		FieldSubject:     fields.NewStringValue(views.EnsureSubject(e.Subject)),
		FieldBodyPreview: fields.NewStringValue(valueOrNotDefined(e.BodyPreview)),
		FieldImportance:  fields.NewStringValue(valueOrNotDefined(e.Importance)),
		FieldWhen:        fields.NewStringValue(valueOrNotDefined(formattedDate)),
		FieldDuration:    fields.NewStringValue(valueOrNotDefined(dur)),
		FieldOrganizer: fields.NewStringValue(
			fmt.Sprintf("[%s](mailto:%s)",
				e.Organizer.EmailAddress.Name, e.Organizer.EmailAddress.Address)),
		FieldLocation:       fields.NewStringValue(valueOrNotDefined(e.Location)),
		FieldResponseStatus: fields.NewStringValue(e.ResponseStatus.Response),
		FieldAttendees:      fields.NewMultiValue(attendees...),
	}

	return ff
}

func valueOrNotDefined(s string) string {
	if s == "" {
		return "Not defined"
	}

	return s
}
