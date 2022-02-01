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

func (c *client) CreateMySubscription(remoteUserEmail string, notificationURL string) (*remote.Subscription, error) {
	// TODO: Use siteURL field from mattermost config
	sub := &remote.Subscription{
		WebhookNotificationUrl: fmt.Sprintf("%s/%s/%s", c.conf.MattermostServerBaseURL, c.conf.PluginURLPath, notificationURL),
	}

	path, err := c.GetEndpointURL(fmt.Sprintf("%s%s", config.PathNotification, config.PathSubscribe), &remoteUserEmail)
	if err != nil {
		return nil, errors.Wrap(err, "ews Subscribe")
	}

	_, err = c.CallJSON(http.MethodPost, path, &sub, &sub)
	if err != nil {
		return nil, errors.Wrap(err, "ews Subscribe")
	}

	c.Logger.With(bot.LogContext{
		"subscriptionID": sub.ID,
	}).Debugf("ews: created subscription.")

	return sub, nil
}

func (c *client) DeleteSubscription(subscriptionID string) error {
	sub := &remote.Subscription{
		ID: subscriptionID,
	}

	path, err := c.GetEndpointURL(fmt.Sprintf("%s%s", config.PathNotification, config.PathUnsubscribe), nil)
	if err != nil {
		return errors.Wrap(err, "ews DeleteSubscription")
	}
	_, err = c.CallJSON(http.MethodPost, path, &sub, nil)
	if err != nil {
		return errors.Wrap(err, "ews DeleteSubscription")
	}

	c.Logger.With(bot.LogContext{
		"subscriptionID": subscriptionID,
	}).Debugf("ews: deleted subscription.")

	return nil
}
