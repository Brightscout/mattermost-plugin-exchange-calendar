package views

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

var prettyStatuses = map[string]string{
	model.StatusOnline:  "Online",
	model.StatusAway:    "Away",
	model.StatusDnd:     "Do Not Disturb",
	model.StatusOffline: "Offline",
}

func RenderStatusChangeNotificationView(events []*remote.Event, status string, customStatus *model.CustomStatus, url string) *model.SlackAttachment {
	nEvents := len(events)
	if nEvents > 0 {
		return statusChangeAttachments(events[0], status, customStatus, url)
	}

	if nEvents > 0 && status == model.StatusDnd {
		return statusChangeAttachments(events[nEvents-1], status, customStatus, url)
	}
	return statusChangeAttachments(nil, status, customStatus, url)
}

func RenderEventWillStartLine(subject, weblink string, startTime time.Time) string {
	eventString := fmt.Sprintf("Your event [%s](%s) will start soon.", subject, weblink)
	if subject == "" {
		eventString = fmt.Sprintf("[An event with no subject](%s) will start soon.", weblink)
	}
	if startTime.Before(time.Now()) {
		eventString = fmt.Sprintf("Your event [%s](%s) is ongoing.", subject, weblink)
		if subject == "" {
			eventString = fmt.Sprintf("[An event with no subject](%s) is ongoing.", weblink)
		}
	}
	return eventString
}

func renderScheduleItem(event *remote.Event, status string) string {
	if event == nil {
		return fmt.Sprintf("You have no upcoming events.\n Shall I change your status back to %s?", prettyStatuses[status])
	}

	resp := RenderEventWillStartLine(event.Subject, event.Weblink, event.Start.Time())

	resp += fmt.Sprintf("\nShall I change your status to %s?", prettyStatuses[status])
	return resp
}

func statusChangeAttachments(event *remote.Event, status string, customStatus *model.CustomStatus, url string) *model.SlackAttachment {
	actionYes := &model.PostAction{
		Name: "Yes",
		Integration: &model.PostActionIntegration{
			URL: url,
			Context: map[string]interface{}{
				"value":              true,
				"change_to":          status,
				"pretty_change_to":   prettyStatuses[status],
				"hasEvent":           false,
				"removeCustomStatus": true,
			},
		},
	}

	actionNo := &model.PostAction{
		Name: "No",
		Integration: &model.PostActionIntegration{
			URL: url,
			Context: map[string]interface{}{
				"value":    false,
				"hasEvent": false,
			},
		},
	}

	if customStatus != nil {
		actionYes.Integration.Context["removeCustomStatus"] = false
		actionYes.Integration.Context["customStatusText"] = customStatus.Text
		actionYes.Integration.Context["customStatusEmoji"] = customStatus.Emoji
		actionYes.Integration.Context["customStatusExpiresAt"] = customStatus.ExpiresAt.String()
		actionYes.Integration.Context["customStatusDuration"] = customStatus.Duration
	}

	if event != nil {
		marshalledStart, _ := json.Marshal(event.Start.Time())
		actionYes.Integration.Context["hasEvent"] = true
		actionYes.Integration.Context["subject"] = event.Subject
		actionYes.Integration.Context["weblink"] = event.Weblink
		actionYes.Integration.Context["startTime"] = string(marshalledStart)

		actionNo.Integration.Context["hasEvent"] = true
		actionNo.Integration.Context["subject"] = event.Subject
		actionNo.Integration.Context["weblink"] = event.Weblink
		actionNo.Integration.Context["startTime"] = string(marshalledStart)
	}

	title := "Status change"
	text := renderScheduleItem(event, status)
	sa := &model.SlackAttachment{
		Title:    title,
		Text:     text,
		Actions:  []*model.PostAction{actionYes, actionNo},
		Fallback: fmt.Sprintf("%s: %s", title, text),
	}

	return sa
}
