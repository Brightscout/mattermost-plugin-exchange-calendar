// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

import (
	"net/url"
	"time"
)

type Client interface {
	AcceptEvent(remoteUserEmail, eventID string) error
	CallFormPost(method, path string, in url.Values, out interface{}) (responseData []byte, err error)
	CallJSON(method, path string, in, out interface{}) (responseData []byte, err error)
	CreateCalendar(remoteUserEmail string, calendar *Calendar) (*Calendar, error)
	CreateEvent(remoteUserEmail string, calendarEvent *Event) (*Event, error)
	CreateMySubscription(remoteUserEmail string, notificationURL string) (*Subscription, error)
	DeclineEvent(remoteUserEmail, eventID string) error
	DeleteCalendar(remoteUserEmail, calendarID string) (*Calendar, error)
	DeleteSubscription(subscriptionID string) error
	FindMeetingTimes(remoteUserEmail string, meetingParams *FindMeetingTimesParameters) (*MeetingTimeSuggestionResults, error)
	GetCalendars(remoteUserEmail string) ([]*Calendar, error)
	GetDefaultCalendarView(remoteUserEmail string, startTime, endTime time.Time) ([]*Event, error)
	DoBatchViewCalendarRequests([]*ViewCalendarParams) ([]*ViewCalendarResponse, error)
	GetEvent(remoteUserEmail, eventID string) (*Event, error)
	GetMailboxSettings(remoteUserID string) (*MailboxSettings, error)
	GetMe(remoteUserEmail string) (*User, error)
	GetNotificationData(remoteUserEmail, eventID, subscriptionID string) (*Notification, error)
	GetSchedule(requests []*ScheduleUserInfo, startTime, endTime *DateTime, availabilityViewInterval int) ([]*ScheduleInformation, error)
	TentativelyAcceptEvent(remoteUserEmail, eventID string) error
	GetSuperuserToken() (string, error)
}
