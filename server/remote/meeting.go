// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

type FindMeetingTimesParameters struct {
	Attendees []Attendee `json:"attendees,omitempty"`
}

type TimeConstraint struct {
	ActivityDomain string     `json:"activityDomain,omitempty"`
	TimeSlots      []TimeSlot `json:"timeSlots,omitempty"`
}
type MeetingTimeSuggestion struct {
	MeetingTimeSlot *DateTime `json:"meetingTimeSlot"`
}

type AttendeeAvailability struct {
	Attendee     *Attendee
	Availability string `json:"availability"`
}

type MeetingTimeSuggestionResults struct {
	MeetingTimeSuggestions []*MeetingTimeSuggestion `json:"meetingTimeSuggestions"`
}

type TimeSlot struct {
	Start *DateTime `json:"start,omitempty"`
	End   *DateTime `json:"end,omitempty"`
}

type LocationConstraint struct {
	Locations       []LocationConstraintItem `json:"locations,omitempty"`
	IsRequired      *bool                    `json:"isRequired,omitempty"`
	SuggestLocation *bool                    `json:"suggestLocation,omitempty"`
}

type LocationConstraintItem struct {
	Location            *Location
	ResolveAvailability *bool `json:"resolveAvailability,omitempty"`
}
