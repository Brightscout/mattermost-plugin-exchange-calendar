package config

import "github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"

// StoredConfig represents the data stored in and managed with the Mattermost
// config.
type StoredConfig struct {
	EWSProxyServerBaseURL string
	EWSProxyServerAuthKey string
	StatusSyncJobInterval int64
	AutoConnectUsers      bool

	EnableStatusSync   bool
	EnableDailySummary bool

	bot.Config
}

// Config represents the the metadata handed to all request runners (command,
// http).
type Config struct {
	StoredConfig

	BuildDate              string
	BuildHash              string
	BuildHashShort         string
	MattermostSiteHostname string
	MattermostSiteURL      string
	PluginID               string
	PluginURL              string
	PluginURLPath          string
	PluginVersion          string
}
