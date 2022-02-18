// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

type FindMeetingTimesParameters struct {
	Attendees []Attendee `json:"attendees,omitempty"`
}

type MeetingTimeSuggestion struct {
	MeetingTimeSlot *DateTime `json:"meetingTimeSlot"`
}

type MeetingTimeSuggestionResults struct {
	MeetingTimeSuggestions []*MeetingTimeSuggestion `json:"meetingTimeSuggestions"`
}
