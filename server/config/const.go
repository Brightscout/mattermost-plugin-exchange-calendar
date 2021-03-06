// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package config

const (
	BotUserName    = "mscalendar"
	BotDisplayName = "Microsoft Calendar"
	BotDescription = "Created by the Microsoft Calendar Plugin."

	ApplicationName    = "Microsoft Calendar"
	Repository         = "mattermost-plugin-exchange-mscalendar"
	CommandTrigger     = "mscalendar"
	TelemetryShortName = "mscalendar"

	PathOAuth2              = "/oauth2"
	PathComplete            = "/complete"
	PathAPI                 = "/api/v1"
	PathPostAction          = "/action"
	PathRespond             = "/respond"
	PathAccept              = "/accept"
	PathDecline             = "/decline"
	PathTentative           = "/tentative"
	PathConfirmStatusChange = "/confirm"
	PathGetNotification     = "/notification"
	PathNotification        = "/api/notification"
	PathEvent               = "/api/event"
	PathCalendar            = "/api/calendar"
	PathFindMeetingTimes    = "/suggestions"
	PathUser                = "/api/user"
	PathSubscribe           = "/subscribe"
	PathUnsubscribe         = "/unsubscribe"
	PathBatch               = "/api/batch"
	PathBatchEvent          = PathBatch + "/event"
	PathBatchUser           = PathBatch + "/user"
	PathBatchSubscription   = PathBatch + PathSubscribe
	PathSync                = "/sync"
	PathSubscription        = "/subscription"

	// TODO: Change path from notification/api/event to /api/notification/event
	FullPathEventNotification = PathGetNotification + PathEvent
	FullPathOAuth2Redirect    = PathOAuth2 + PathComplete

	EventIDKey = "EventID"
	EmailKey   = "email"

	CreateCalendarHeading = "Calendar created."
	DeleteCalendarHeading = "Calendar deleted."

	AuthorizationHeaderKey = "Authorization"
	UsersCountPerPage      = 20

	Organizer = "organizer"
	Required  = "required"
	Optional  = "optional"

	CustomStatusEmoji = "calendar"
	CustomStatusText  = "In a meeting"

	SubscriptionStatusOK          = "OK"
	SubscriptionStatusUnsubscribe = "Unsubscribe"
)
