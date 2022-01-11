// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package mscalendar

import (
	"context"
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/oauth2connect"
)

const BotWelcomeMessage = "Bot user connected to account %s."

const RemoteUserAlreadyConnected = "%s account `%s` is already mapped to Mattermost account `%s`. Please run `/%s disconnect`, while logged in as the Mattermost account"
const RemoteUserAlreadyConnectedNotFound = "%s account `%s` is already mapped to a Mattermost account, but the Mattermost user could not be found"

type oauth2App struct {
	Env
}

func NewOAuth2App(env Env) oauth2connect.App {
	return &oauth2App{
		Env: env,
	}
}

func (app *oauth2App) InitOAuth2(mattermostUserID string) (url string, err error) {
	user, err := app.Store.LoadUser(mattermostUserID)
	if err == nil {
		return "", fmt.Errorf("user is already connected to %s", user.Remote.Mail)
	}

	conf := app.Remote.NewOAuth2Config()
	state := fmt.Sprintf("%v_%v", model.NewId()[0:15], mattermostUserID)
	err = app.Store.StoreOAuth2State(state)
	if err != nil {
		return "", err
	}

	return conf.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (app *oauth2App) CompleteOAuth2(authedUserID string) error {
	if authedUserID == "" {
		return errors.New("missing user")
	}

	// oconf := app.Remote.NewOAuth2Config()

	// err := app.Store.VerifyOAuth2State(state)
	// if err != nil {
	// 	return errors.WithMessage(err, "missing stored state")
	// }

	// mattermostUserID := strings.Split(state, "_")[1]
	// if mattermostUserID != authedUserID {
	// 	return errors.New("not authorized, user ID mismatch")
	// }

	ctx := context.Background()
	// tok, err := oconf.Exchange(ctx, code)
	// if err != nil {
	// 	return err
	// }

	client := app.Remote.MakeClient(ctx)
	user, userErr := app.PluginAPI.GetMattermostUser(authedUserID)
	if userErr != nil {
		return userErr
	}

	me, err := client.GetMe(user.Email)
	if err != nil {
		return err
	}

	_, err = app.Store.LoadMattermostUserID(me.ID)
	if err == nil {
		// Couldn't fetch connected MM account. Reject connect attempt.
		app.Poster.DM(authedUserID, RemoteUserAlreadyConnectedNotFound, config.ApplicationName, me.Mail)
		return fmt.Errorf(RemoteUserAlreadyConnectedNotFound, config.ApplicationName, me.Mail)
	}

	u := &store.User{
		PluginVersion:    app.Config.PluginVersion,
		MattermostUserID: authedUserID,
		Remote:           me,
	}

	mailboxSettings, err := client.GetMailboxSettings(me.ID)
	if err != nil {
		return err
	}

	u.Settings.DailySummary = &store.DailySummaryUserSettings{
		PostTime: "8:00AM",
		Timezone: mailboxSettings.TimeZone,
		Enable:   false,
	}

	err = app.Store.StoreUser(u)
	if err != nil {
		return err
	}

	err = app.Store.StoreUserInIndex(u)
	if err != nil {
		return err
	}

	app.Welcomer.AfterSuccessfullyConnect(authedUserID, me.Mail)

	return nil
}
