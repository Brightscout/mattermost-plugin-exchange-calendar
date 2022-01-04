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

func (c *client) DeleteCalendar(remoteUserEmail string, calID string) error {
	calOut := &remote.Calendar{}
	url, err := c.GetEndpointURL(remoteUserEmail, fmt.Sprintf("%s/%s", config.PathCalendar, calID))
	fmt.Println("URL is ", url)
	if err != nil {
		return errors.Wrap(err, "msgraph DeleteCalendar")
	}
	_, err = c.CallJSON(http.MethodDelete, url, nil, calOut)
	if err != nil {
		return errors.Wrap(err, "msgraph DeleteCalendar")
	}
	c.Logger.With(bot.LogContext{}).Infof("msgraph: DeleteCalendar deleted calendar `%v`.", calID)
	return nil
}
