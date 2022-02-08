// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
)

func (c *client) GetNotificationData(remoteUserEmail, eventID, subscriptionID string) (*remote.Notification, error) {
	event := &remote.Event{}
	path, err := c.GetEndpointURL(fmt.Sprintf("%s/%s", config.PathEvent, eventID), &remoteUserEmail)
	if err != nil {
		return nil, errors.Wrap(err, "ews GetNotificationData")
	}
	_, err = c.CallJSON(http.MethodGet, path, nil, &event)
	if err != nil {
		c.Logger.With(bot.LogContext{
			"subscriptionID": subscriptionID,
		}).Infof("ews: failed to fetch notification data resource: `%v`.", err)
		return nil, errors.Wrap(err, "ews GetNotificationData")
	}

	notification := &remote.Notification{
		Event: event,
	}
	return notification, nil
}
