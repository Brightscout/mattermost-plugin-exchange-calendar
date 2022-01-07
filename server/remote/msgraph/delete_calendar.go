// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"fmt"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
	"github.com/pkg/errors"
)

func (c *client) DeleteCalendar(remoteUserEmail string, calID string) (*remote.Calendar, error) {
	calOut := &remote.Calendar{}
	url, err := c.GetEndpointURL(remoteUserEmail, fmt.Sprintf("%s/%s", config.PathCalendar, calID))
	if err != nil {
		return nil, errors.Wrap(err, "ews DeleteCalendar")
	}
	_, err = c.CallJSON(http.MethodDelete, url, nil, calOut)
	if err != nil {
		return nil, errors.Wrap(err, "ews DeleteCalendar")
	}
	c.Logger.With(bot.LogContext{}).Infof("ews: DeleteCalendar deleted calendar `%v`.", calID)
	return calOut, nil
}
