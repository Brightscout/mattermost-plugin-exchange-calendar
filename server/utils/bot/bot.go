// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package bot

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/flow"
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/pkg/errors"
)

type Bot interface {
	Poster
	Logger
	Admin
	FlowController

	Ensure(stored *model.Bot, iconPath string) error
	WithConfig(Config) Bot
	MattermostUserID() string
	RegisterFlow(flow.Flow, flow.Store)
}

type bot struct {
	Config
	pluginAPI        plugin.API
	client           pluginapi.Client
	mattermostUserID string
	displayName      string
	logContext       LogContext
	pluginURL        string

	flow      flow.Flow
	flowStore flow.Store
}

func New(api plugin.API, client pluginapi.Client, pluginURL string) Bot {
	return &bot{
		pluginAPI: api,
		client:    client,
		pluginURL: pluginURL,
	}
}

func (bot *bot) RegisterFlow(flow flow.Flow, flowStore flow.Store) {
	bot.flow = flow
	bot.flowStore = flowStore
}

func (bot *bot) Ensure(stored *model.Bot, iconPath string) error {
	if bot.mattermostUserID != "" {
		// Already done
		return nil
	}

	botUserID, err := bot.client.Bot.EnsureBot(stored)
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot account")
	}
	bot.mattermostUserID = botUserID
	bot.displayName = stored.DisplayName
	return nil
}

func (bot *bot) WithConfig(conf Config) Bot {
	newbot := *bot
	newbot.Config = conf
	return &newbot
}

func (bot *bot) MattermostUserID() string {
	return bot.mattermostUserID
}

func (bot *bot) String() string {
	return bot.displayName
}
