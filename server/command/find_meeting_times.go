// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
)

func (c *Command) findMeetings(parameters ...string) (string, bool, error) {
	meetingParams := &remote.FindMeetingTimesParameters{}

	var attendees []remote.Attendee
	for a := range parameters {
		s := strings.Split(parameters[a], ":")
		t, email := s[0], s[1]
		attendee := remote.Attendee{
			Type: t,
			EmailAddress: &remote.EmailAddress{
				Address: email,
			},
		}
		attendees = append(attendees, attendee)
	}
	meetingParams.Attendees = attendees

	meetings, err := c.MSCalendar.FindMeetingTimes(c.user(), meetingParams)
	if err != nil {
		return "", false, err
	}

	timeZone, _ := c.MSCalendar.GetTimezone(c.user())
	resp := ""
	for _, m := range meetings.MeetingTimeSuggestions {
		if timeZone != "" {
			m.MeetingTimeSlot = m.MeetingTimeSlot.In(timeZone)
			m.MeetingTimeSlot.TimeZone = timeZone
		}
	}

	sort.Slice(meetings.MeetingTimeSuggestions, func(i, j int) bool {
		return meetings.MeetingTimeSuggestions[i].MeetingTimeSlot.Time().Before(meetings.MeetingTimeSuggestions[j].MeetingTimeSlot.Time())
	})

	for _, m := range meetings.MeetingTimeSuggestions {
		resp += utils.JSONBlock(renderMeetingTime(m))
	}

	return resp, false, nil
}

func renderMeetingTime(m *remote.MeetingTimeSuggestion) string {
	return fmt.Sprintf("%s (%s)", m.MeetingTimeSlot.PrettyString(), m.MeetingTimeSlot.TimeZone)
}
