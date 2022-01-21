// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package mscalendar

import (
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
)

type Subscriptions interface {
	CreateMyEventSubscription() (*store.Subscription, error)
	DeleteMyEventSubscription() error
	LoadMyEventSubscription() (*store.Subscription, error)
}

func (m *mscalendar) CreateMyEventSubscription() (*store.Subscription, error) {
	err := m.Filter(withClient)
	if err != nil {
		return nil, err
	}

	sub, err := m.client.CreateMySubscription(m.actingUser.Remote.Mail, config.FullPathEventNotification)
	if err != nil {
		return nil, err
	}

	storedSub := &store.Subscription{
		Remote:              sub,
		MattermostCreatorID: m.actingUser.MattermostUserID,
		PluginVersion:       m.Config.PluginVersion,
	}
	err = m.Store.StoreUserSubscription(m.actingUser.User, storedSub)
	if err != nil {
		return nil, err
	}

	return storedSub, nil
}

func (m *mscalendar) LoadMyEventSubscription() (*store.Subscription, error) {
	err := m.Filter(withActingUserExpanded)
	if err != nil {
		return nil, err
	}
	storedSub, err := m.Store.LoadSubscription(m.actingUser.Settings.EventSubscriptionID)
	if err != nil {
		return nil, err
	}
	return storedSub, err
}

func (m *mscalendar) DeleteMyEventSubscription() error {
	err := m.Filter(withActingUserExpanded)
	if err != nil {
		return err
	}

	subscriptionID := m.actingUser.Settings.EventSubscriptionID

	err = m.client.DeleteSubscription(subscriptionID)
	if err != nil {
		return errors.WithMessagef(err, "failed to delete subscription %s", subscriptionID)
	}

	err = m.Store.DeleteUserSubscription(m.actingUser.User, subscriptionID)
	if err != nil {
		return errors.WithMessagef(err, "failed to delete subscription %s", subscriptionID)
	}

	return nil
}
