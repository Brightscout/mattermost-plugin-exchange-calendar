package views

import (
	"encoding/json"
	"fmt"
	URL "net/url"
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

func RenderStatusChangeNotificationView(events []*remote.Event, status, url string) *model.SlackAttachment {
	fmt.Print("inside RenderStatusChangeNotificationView\n")
	for _, e := range events {
		if e.Start.Time().After(time.Now()) {
			return statusChangeAttachments(e, status, url)
		}
	}

	nEvents := len(events)
	if nEvents > 0 && status == model.StatusDnd {
		return statusChangeAttachments(events[nEvents-1], status, url)
	}
	return statusChangeAttachments(nil, status, url)
}

func RenderCustomStatusChangeNotificationView(events []*remote.Event, url string) *model.SlackAttachment {
	fmt.Print("inside RenderCustomStatusChangeNotificationView Custom\n")

	for _, e := range events {
		if e.Start.Time().After(time.Now()) {
			return customstatusChangeAttachments(e, url)
		}
	}
	return customstatusChangeAttachments(nil, url)
}

func RenderEventWillStartLine(subject, weblink string, startTime time.Time) string {
	link, _ := URL.QueryUnescape(weblink)
	eventString := fmt.Sprintf("Your event [%s](%s) will start soon.", subject, link)
	if subject == "" {
		eventString = fmt.Sprintf("[An event with no subject](%s) will start soon.", link)
	}
	if startTime.Before(time.Now()) {
		eventString = fmt.Sprintf("Your event [%s](%s) is ongoing.", subject, link)
		if subject == "" {
			eventString = fmt.Sprintf("[An event with no subject](%s) is ongoing.", link)
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

func statusChangeAttachments(event *remote.Event, status, url string) *model.SlackAttachment {
	actionYes := &model.PostAction{
		Name: "Yes",
		Integration: &model.PostActionIntegration{
			URL: url,
			Context: map[string]interface{}{
				"value":            true,
				"change_to":        status,
				"pretty_change_to": prettyStatuses[status],
				"hasEvent":         false,
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

func customstatusChangeAttachments(event *remote.Event, url string) *model.SlackAttachment {
	actionYes := &model.PostAction{
		Name: "Yes",
		Integration: &model.PostActionIntegration{
			URL: url,
			Context: map[string]interface{}{
				"value": true,
				"hasEvent": false,
			},
		},
	}

	actionNo := &model.PostAction{
		Name: "No",
		Integration: &model.PostActionIntegration{
			URL: url,
			Context: map[string]interface{}{
				"value": false,
				"hasEvent": false,
			},
		},
	}

	if event != nil {
		marshalledStart, _ := json.Marshal(event.Start.Time())

		actionYes.Integration.Context["hasEvent"] = true
		actionYes.Integration.Context["subject"] = event.Subject
		actionYes.Integration.Context["weblink"] = event.Weblink
		actionYes.Integration.Context["startTime"] = string(marshalledStart)
		actionYes.Integration.Context["endTime"] = event.End.String()

		actionNo.Integration.Context["hasEvent"] = true
		actionNo.Integration.Context["subject"] = event.Subject
		actionNo.Integration.Context["weblink"] = event.Weblink
		actionNo.Integration.Context["startTime"] = string(marshalledStart)
		actionNo.Integration.Context["endTime"] = event.End.String()
	}
	title := "Custom Status change"
	//link, _ := URL.QueryUnescape(event.Weblink)
	//text := fmt.Sprintf("Your event [%s](%s) will start soon.\n", event.Subject, link)
	var sa *model.SlackAttachment
	if(event==nil){
		actionYes.Integration.Context["setStatus"] = false
		actionNo.Integration.Context["setStatus"] = false
		sa = &model.SlackAttachment{
			Title:    title,
			Text:     "You have no upcoming events \nDo you want to unset your custom status",
			Actions:  []*model.PostAction{actionYes, actionNo},
			Fallback: fmt.Sprintf("%s: %s", title, "Do you want to unset your custom status"),
		}
	} else{
		actionYes.Integration.Context["setStatus"] = true
		actionNo.Integration.Context["setStatus"] = true
		sa = &model.SlackAttachment{
		Title:    title,
		Text:     "Do you want to change your custom status to In a Meeting",
		Actions:  []*model.PostAction{actionYes, actionNo},
		Fallback: fmt.Sprintf("%s: %s", title, "Do you want to change your custom status to In a Meeting"),
	}
}
	return sa
}
