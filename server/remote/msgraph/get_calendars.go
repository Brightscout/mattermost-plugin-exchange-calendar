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

func (c *client) GetCalendars(remoteUserEmail string) ([]*remote.Calendar, error) {
	calOut := []*remote.Calendar{}
	url, err := c.GetEndpointURL(config.PathCalendar, &remoteUserEmail)
	if err != nil {
		return nil, errors.Wrap(err, "msgraph GetCalendars")
	}
	_, err = c.CallJSON(http.MethodGet, url, nil, &calOut)
	if err != nil {
		return nil, errors.Wrap(err, "msgraph GetCalendars")
	}
	c.Logger.With(bot.LogContext{
		"UserID": remoteUserEmail,
		"v":      calOut,
	}).Infof("msgraph: GetUserCalendars returned `%d` calendars.", len(calOut))
	return calOut, nil
}
