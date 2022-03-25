// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"encoding/json"
	"fmt"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/kvstore"
)

type SubscriptionStore interface {
	LoadSubscription(subscriptionID string) (*Subscription, error)
	LoadSubscriptionIndex() (SubscriptionIndex, error)
	StoreUserSubscription(user *User, subscription *Subscription) error
	DeleteUserSubscription(user *User, subscriptionID string) error
	DeleteSubscriptionIndex() error
	StoreUserSubscriptionInIndex(user *User, subscription *Subscription) error
	DeleteUserSubscriptionFromIndex(subscriptionID string) error
}

type SubscriptionIndex []*SubscriptionShort

type SubscriptionShort struct {
	MattermostCreatorID string `json:"mm_id"`
	Email               string `json:"email"`
	SubscriptionID      string `json:"subscription_id"`
}

type Subscription struct {
	PluginVersion       string
	Remote              *remote.Subscription
	MattermostCreatorID string
}

func (s *pluginStore) LoadSubscription(subscriptionID string) (*Subscription, error) {
	sub := Subscription{}
	err := kvstore.LoadJSON(s.subscriptionKV, subscriptionID, &sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *pluginStore) LoadSubscriptionIndex() (SubscriptionIndex, error) {
	subs := SubscriptionIndex{}
	err := kvstore.LoadJSON(s.subscriptionIndexKV, "", &subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *pluginStore) DeleteSubscriptionIndex() error {
	err := s.subscriptionIndexKV.Delete("")
	if err != nil {
		return err
	}
	return nil
}

func (s *pluginStore) StoreUserSubscription(user *User, subscription *Subscription) error {
	if user.Remote.ID != subscription.Remote.CreatorID {
		return fmt.Errorf("user %q does not match the subscription creator %q",
			user.Remote.ID, subscription.Remote.CreatorID)
	}
	err := kvstore.StoreJSON(s.subscriptionKV, subscription.Remote.ID, subscription)
	if err != nil {
		return err
	}
	user.Settings.EventSubscriptionID = subscription.Remote.ID
	err = kvstore.StoreJSON(s.userKV, user.MattermostUserID, user)
	if err != nil {
		return err
	}

	s.Logger.With(bot.LogContext{
		"mattermostUserID": user.MattermostUserID,
		"remoteUserID":     subscription.Remote.CreatorID,
		"subscriptionID":   subscription.Remote.ID,
	}).Debugf("store: stored mattermost user subscription.")
	return nil
}

func (s *pluginStore) DeleteUserSubscription(user *User, subscriptionID string) error {
	err := s.subscriptionKV.Delete(subscriptionID)
	if err != nil {
		return err
	}
	mattermostUserID := ""
	if user != nil {
		user.Settings.EventSubscriptionID = ""
		err = s.StoreUser(user)
		if err != nil {
			return err
		}
		mattermostUserID = user.MattermostUserID
	}

	var subscriptionIndex []*SubscriptionShort
	err = kvstore.LoadJSON(s.subscriptionIndexKV, "", &subscriptionIndex)
	if err != nil {
		return err
	}
	filtered := []*SubscriptionShort{}
	for _, s := range subscriptionIndex {
		if s.SubscriptionID != subscriptionID {
			filtered = append(filtered, s)
		}
	}
	err = kvstore.StoreJSON(s.subscriptionIndexKV, "", &filtered)
	if err != nil {
		return err
	}

	s.Logger.With(bot.LogContext{
		"mattermostUserID": mattermostUserID,
		"subscriptionID":   subscriptionID,
	}).Debugf("store: deleted mattermost user subscription.")
	return nil
}

func (s *pluginStore) ModifySubscriptionIndex(modify func(subscriptionIndex SubscriptionIndex) (SubscriptionIndex, error)) error {
	return kvstore.AtomicModify(s.subscriptionIndexKV, "", func(initial []byte, storeErr error) ([]byte, error) {
		if storeErr != nil && storeErr != ErrNotFound {
			return initial, storeErr
		}

		var storedIndex SubscriptionIndex
		if len(initial) > 0 {
			err := json.Unmarshal(initial, &storedIndex)
			if err != nil {
				return nil, err
			}
		}

		updated, err := modify(storedIndex)
		if err != nil {
			return nil, err
		}

		result, err := json.Marshal(updated)
		if err != nil {
			return nil, err
		}

		return result, nil
	})
}

func (s *pluginStore) StoreUserSubscriptionInIndex(user *User, subscription *Subscription) error {
	return s.ModifySubscriptionIndex(func(subscriptionIndex SubscriptionIndex) (SubscriptionIndex, error) {
		newSub := &SubscriptionShort{
			MattermostCreatorID: subscription.MattermostCreatorID,
			SubscriptionID:      subscription.Remote.ID,
			Email:               user.Remote.ID,
		}

		for i, s := range subscriptionIndex {
			if s.MattermostCreatorID == subscription.MattermostCreatorID && s.SubscriptionID == subscription.Remote.ID {
				// Removing old subscription of user and adding the new subscription in array
				result := append(subscriptionIndex[:i], newSub)
				return append(result, subscriptionIndex[i+1:]...), nil
			}
		}

		return append(subscriptionIndex, newSub), nil
	})
}

func (s *pluginStore) DeleteUserSubscriptionFromIndex(subscriptionID string) error {
	return s.ModifySubscriptionIndex(func(subscriptionIndex SubscriptionIndex) (SubscriptionIndex, error) {
		for i, s := range subscriptionIndex {
			if s.SubscriptionID == subscriptionID {
				// Removing subscription data of user from array
				return append(subscriptionIndex[:i], subscriptionIndex[i+1:]...), nil
			}
		}
		return subscriptionIndex, nil
	})
}
