// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"errors"
	"fmt"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
)

const (
	ConnectBotAlreadyConnectedTemplate = "The bot account is already connected to %s account `%s`. To connect to a different account, first run `/%s disconnect_bot`."
	ConnectBotSuccessTemplate          = "[Click here to link the bot's %s account.](%s/oauth2/connect_bot)"
	ConnectAlreadyConnectedTemplate    = "Your Mattermost account is already connected to %s account `%s`. To connect to a different account, first run `/%s disconnect`."
	ConnectErrorMessage                = "There has been a problem while trying to connect. err="
)

func (c *Command) connect(parameters ...string) (string, bool, error) {
	ru, err := c.MSCalendar.GetRemoteUser(c.Args.UserId)
	if err == nil {
		return fmt.Sprintf(ConnectAlreadyConnectedTemplate, config.ApplicationName, ru.Mail, config.CommandTrigger), false, nil
	}

	err = c.MSCalendar.CompleteOAuth2(c.Args.UserId)
	if err != nil {
		return "", false, errors.New(ConnectErrorMessage + err.Error())
	}

	return "", false, nil
}
