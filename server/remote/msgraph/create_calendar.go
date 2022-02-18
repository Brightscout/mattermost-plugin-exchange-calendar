// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
	"github.com/pkg/errors"
)

// CreateCalendar creates a calendar
func (c *client) CreateCalendar(remoteUserEmail string, calIn *remote.Calendar) (*remote.Calendar, error) {
	calOut := &remote.Calendar{}
	url, err := c.GetEndpointURL(config.PathCalendar, &remoteUserEmail)
	if err != nil {
		return nil, errors.Wrap(err, "ews CreateCalendar")
	}
	_, err = c.CallJSON(http.MethodPost, url, calIn, calOut)
	if err != nil {
		return nil, errors.Wrap(err, "ews CreateCalendar")
	}
	c.Logger.With(bot.LogContext{
		"v": calOut,
	}).Infof("ews: CreateCalendar created the following calendar.")
	return calOut, nil
}
