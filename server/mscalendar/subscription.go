// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package mscalendar

import (
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
)

type Subscriptions interface {
	CreateMyEventSubscription() (*store.Subscription, error)
	DeleteOrphanedSubscription(ID string) error
	DeleteMyEventSubscription() error
	LoadMyEventSubscription() (*store.Subscription, error)
	SyncUserSubscriptions() error
	GetSubscriptionByID(subscriptionID string) (*store.Subscription, error)
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
	err = m.Store.StoreUserSubscriptionInIndex(m.actingUser.User, storedSub)
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

	err = m.DeleteOrphanedSubscription(subscriptionID)
	if err != nil {
		return err
	}

	err = m.Store.DeleteUserSubscription(m.actingUser.User, subscriptionID)
	if err != nil {
		return errors.WithMessagef(err, "failed to delete subscription %s", subscriptionID)
	}

	return nil
}

func (m *mscalendar) GetSubscriptionByID(subscriptionID string) (*store.Subscription, error) {
	storedSub, err := m.Store.LoadSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	return storedSub, err
}

func (m *mscalendar) DeleteOrphanedSubscription(subscriptionID string) error {
	err := m.Filter(withClient)
	if err != nil {
		return err
	}
	err = m.client.DeleteSubscription(m.actingUser.Remote.Mail, subscriptionID)
	if err != nil {
		return errors.WithMessagef(err, "failed to delete subscription %s", subscriptionID)
	}
	return nil
}

func (m *mscalendar) SyncUserSubscriptions() error {
	// Load all user subscriptions
	subscriptionIndex, err := m.Store.LoadSubscriptionIndex()
	if err != nil {
		return err
	}
	totalSubscribedUsers := len(subscriptionIndex)
	if totalSubscribedUsers == 0 {
		return nil
	}

	err = m.Filter(withSuperuserClient)
	if err != nil {
		return err
	}
	// Create new subscription for users and store it
	var requests []remote.SubscriptionBatchSingleRequest
	emailUserMap := make(map[string]*store.User)
	for _, sub := range subscriptionIndex {
		user, err := m.Store.LoadUser(sub.MattermostCreatorID)
		if err != nil {
			m.Logger.Warnf("Error occurred while fetching user from store with id %s. err=%s", sub.MattermostCreatorID, err.Error())
			continue
		}
		// Deleting previous subscription for this user from subscription store
		err = m.Store.DeleteUserSubscription(user, sub.SubscriptionID)
		if err != nil {
			m.Logger.Warnf("failed to delete subscription %s. err=%s", sub.SubscriptionID, err.Error())
			continue
		}
		request := remote.SubscriptionBatchSingleRequest{
			Email: sub.Email,
			Subscription: remote.Subscription{
				WebhookNotificationUrl: m.client.GetWebhookNotificationURL(),
			},
		}
		requests = append(requests, request)
		emailUserMap[sub.Email] = user
	}
	responses, err := m.client.DoBatchSubscriptionRequests(requests)
	if err != nil {
		return err
	}
	for _, response := range responses {
		if response.Error != nil {
			m.Logger.Warnf("Error occurred while subscribing user with email %s. err=%s", response.Email, response.Error.Message)
			continue
		}
		user, exists := emailUserMap[response.Email]
		if !exists {
			m.Logger.Warnf("Error occurred while fetching user with email %s. err=%s", response.Email, response.Error.Message)
			continue
		}
		response.Subscription.CreatorID = response.Email
		storedSub := &store.Subscription{
			Remote:              response.Subscription,
			MattermostCreatorID: user.MattermostUserID,
			PluginVersion:       m.Config.PluginVersion,
		}
		err = m.Store.StoreUserSubscription(user, storedSub)
		if err != nil {
			return err
		}
		err = m.Store.StoreUserSubscriptionInIndex(user, storedSub)
		if err != nil {
			return err
		}
	}

	return nil
}
