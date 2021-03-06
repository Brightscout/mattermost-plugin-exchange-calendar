// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
)

type Remote interface {
	MakeClient(context.Context) Client
	MakeSuperuserClient(ctx context.Context) (Client, error)
	NewOAuth2Config() *oauth2.Config
	HandleWebhook(http.ResponseWriter, *http.Request) (bool, *Notification, error)
}

var Makers = map[string]func(*config.Config, bot.Logger) Remote{}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
