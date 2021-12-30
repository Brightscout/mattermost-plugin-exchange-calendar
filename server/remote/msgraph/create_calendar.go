// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
)

// CreateCalendar creates a calendar
func (c *client) CreateCalendar(remoteUserID string, calIn *remote.Calendar) (*remote.Calendar, error) {
	var calOut = remote.Calendar{}
	// TODO: ADD calendar API
	// err := c.rbuilder.Users().ID(remoteUserID).Calendars().Request().JSONRequest(c.ctx, http.MethodPost, "", &calIn, &calOut)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "msgraph CreateCalendar")
	// }
	c.Logger.With(bot.LogContext{
		"v": calOut,
	}).Infof("msgraph: CreateCalendar created the following calendar.")
	return &calOut, nil
}
