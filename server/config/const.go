// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package config

const (
	BotUserName    = "mscalendar"
	BotDisplayName = "Microsoft Calendar"
	BotDescription = "Created by the Microsoft Calendar Plugin."

	ApplicationName    = "Microsoft Calendar"
	Repository         = "mattermost-plugin-mscalendar"
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
	PathMe                  = "/api/me"
	PathSubscribe           = "/subscribe"
	PathUnsubscribe         = "/unsubscribe"
	PathBatch               = "/api/batch"
	PathBatchEvent          = PathBatch + "/event"

	// TODO: Change path from notification/api/event to /api/notification/event
	FullPathEventNotification = PathGetNotification + PathEvent
	FullPathOAuth2Redirect    = PathOAuth2 + PathComplete

	EventIDKey = "EventID"
	EmailKey   = "email"

	CreateCalendarHeading = "Calendar created."
	DeleteCalendarHeading = "Calendar deleted."

	AuthorizationHeaderKey = "Authorization"
)
