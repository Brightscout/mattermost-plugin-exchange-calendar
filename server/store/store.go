// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"time"

	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/tracker"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/flow"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/kvstore"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/settingspanel"
)

const (
	UserKeyPrefix              = "user_"
	UserIndexKeyPrefix         = "userindex_"
	MattermostUserIDKeyPrefix  = "mmuid_"
	OAuth2KeyPrefix            = "oauth2_"
	SubscriptionKeyPrefix      = "sub_"
	SubscriptionIndexKeyPrefix = "subindex_"
	EventKeyPrefix             = "ev_"
	WelcomeKeyPrefix           = "welcome_"
	SettingsPanelPrefix        = "settings_panel_"
)

const OAuth2KeyExpiration = 15 * time.Minute

var ErrNotFound = kvstore.ErrNotFound

type Store interface {
	UserStore
	OAuth2StateStore
	SubscriptionStore
	EventStore
	WelcomeStore
	flow.Store
	settingspanel.SettingStore
	settingspanel.PanelStore
}

type pluginStore struct {
	basicKV             kvstore.KVStore
	oauth2KV            kvstore.KVStore
	userKV              kvstore.KVStore
	mattermostUserIDKV  kvstore.KVStore
	userIndexKV         kvstore.KVStore
	subscriptionKV      kvstore.KVStore
	subscriptionIndexKV kvstore.KVStore
	eventKV             kvstore.KVStore
	welcomeIndexKV      kvstore.KVStore
	settingsPanelKV     kvstore.KVStore
	Logger              bot.Logger
	Tracker             tracker.Tracker
}

func NewPluginStore(api plugin.API, logger bot.Logger, tracker tracker.Tracker) Store {
	basicKV := kvstore.NewPluginStore(api)
	return &pluginStore{
		basicKV:             basicKV,
		userKV:              kvstore.NewHashedKeyStore(basicKV, UserKeyPrefix),
		userIndexKV:         kvstore.NewHashedKeyStore(basicKV, UserIndexKeyPrefix),
		mattermostUserIDKV:  kvstore.NewHashedKeyStore(basicKV, MattermostUserIDKeyPrefix),
		subscriptionKV:      kvstore.NewHashedKeyStore(basicKV, SubscriptionKeyPrefix),
		subscriptionIndexKV: kvstore.NewHashedKeyStore(basicKV, SubscriptionIndexKeyPrefix),
		eventKV:             kvstore.NewHashedKeyStore(basicKV, EventKeyPrefix),
		oauth2KV:            kvstore.NewHashedKeyStore(kvstore.NewOneTimePluginStore(api, OAuth2KeyExpiration), OAuth2KeyPrefix),
		welcomeIndexKV:      kvstore.NewHashedKeyStore(basicKV, WelcomeKeyPrefix),
		settingsPanelKV:     kvstore.NewHashedKeyStore(basicKV, SettingsPanelPrefix),
		Logger:              logger,
		Tracker:             tracker,
	}
}
