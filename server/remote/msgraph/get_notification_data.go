// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"fmt"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
	"github.com/pkg/errors"
)

func (c *client) GetNotificationData(remoteUserEmail string, eventId string) (*remote.Notification, error) {
	var event = remote.Event{}
	var notification = &remote.Notification{}

	path, err := c.GetEndpointURL(remoteUserEmail, fmt.Sprintf("%s/%s", "/api/event", eventId))
	if err != nil {
		return nil, errors.Wrap(err, "ews GetNotificationData")
	}
	_, err = c.CallJSON(http.MethodGet, path, nil, &event)
	if err != nil {
		c.Logger.With(bot.LogContext{
			"subscriptionID": notification.SubscriptionID,
		}).Infof("ews: failed to fetch notification data resource: `%v`.", err)
		return nil, errors.Wrap(err, "ews GetNotificationData")
	}

	notification = &remote.Notification{
		Event: &event,
	}
	return notification, nil
}
