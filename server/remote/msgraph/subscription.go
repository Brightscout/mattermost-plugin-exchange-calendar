// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
)

const subscribeTTL = 48 * time.Hour

func newRandomString() string {
	b := make([]byte, 96)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (c *client) CreateMySubscription(notificationURL string) (*remote.Subscription, error) {
	sub := &remote.Subscription{
		Resource:           "me/events",
		ChangeType:         "created,updated,deleted",
		NotificationURL:    notificationURL,
		ExpirationDateTime: time.Now().Add(subscribeTTL).Format(time.RFC3339),
		ClientState:        newRandomString(),
	}
	// TODO: Add subscription API
	// err := c.rbuilder.Subscriptions().Request().JSONRequest(c.ctx, http.MethodPost, "", sub, sub)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "msgraph CreateMySubscription")
	// }

	c.Logger.With(bot.LogContext{
		"subscriptionID":     sub.ID,
		"resource":           sub.Resource,
		"changeType":         sub.ChangeType,
		"expirationDateTime": sub.ExpirationDateTime,
	}).Debugf("msgraph: created subscription.")

	return sub, nil
}

func (c *client) DeleteSubscription(subscriptionID string) error {
	// TODO: Add subscription API
	// err := c.rbuilder.Subscriptions().ID(subscriptionID).Request().Delete(c.ctx)
	// if err != nil {
	// 	return errors.Wrap(err, "msgraph DeleteSubscription")
	// }

	c.Logger.With(bot.LogContext{
		"subscriptionID": subscriptionID,
	}).Debugf("msgraph: deleted subscription.")

	return nil
}

func (c *client) RenewSubscription(subscriptionID string) (*remote.Subscription, error) {
	sub := remote.Subscription{}
	expires := time.Now().Add(subscribeTTL)
	// TODO: Add subscription API
	// v := struct {
	// 	ExpirationDateTime string `json:"expirationDateTime"`
	// }{
	// 	expires.Format(time.RFC3339),
	// }
	// err := c.rbuilder.Subscriptions().ID(subscriptionID).Request().JSONRequest(c.ctx, http.MethodPatch, "", v, &sub)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "msgraph RenewSubscription")
	// }

	c.Logger.With(bot.LogContext{
		"subscriptionID":     subscriptionID,
		"expirationDateTime": expires.Format(time.RFC3339),
	}).Debugf("msgraph: renewed subscription.")

	return &sub, nil
}

func (c *client) ListSubscriptions() ([]*remote.Subscription, error) {
	var v struct {
		Value []*remote.Subscription `json:"value"`
	}
	// TODO: Add subscription API
	// err := c.rbuilder.Subscriptions().Request().JSONRequest(c.ctx, http.MethodGet, "", nil, &v)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "msgraph ListSubscriptions")
	// }
	return v.Value, nil
}
