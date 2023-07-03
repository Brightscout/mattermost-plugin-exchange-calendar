// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package mscalendar

import (
	"context"
	"fmt"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
)

const BotWelcomeMessage = "Bot user connected to account %s."

const RemoteUserAlreadyConnected = "%s account `%s` is already mapped to Mattermost account `%s`. Please run `/%s disconnect`, while logged in as the Mattermost account"
const RemoteUserAlreadyConnectedNotFound = "%s account `%s` is already mapped to a Mattermost account, but the Mattermost user could not be found"

type OAuth2 interface {
	CompleteOAuth2ForUsers(users []*model.User) error
	CompleteOAuth2(mattermostUserID string) error
	CompleteUserConnect(authedUserID string, timzone model.StringMap, me *remote.User) error
}

func (m *mscalendar) CompleteOAuth2ForUsers(users []*model.User) error {
	var emails []*string
	emailUserMap := make(map[string]*model.User)
	for _, user := range users {
		_, err := m.GetRemoteUser(user.Id)
		if err == nil {
			// User already connected to ms-calendar
			continue
		}
		emails = append(emails, &user.Email)
		emailUserMap[user.Email] = user
	}

	if len(emails) == 0 {
		return nil
	}

	ctx := context.Background()
	client := m.Remote.MakeClient(ctx)
	usersDetails, err := client.GetUsers(emails)
	if err != nil {
		return err
	}

	for _, userDetails := range usersDetails {
		if userDetails.Error != nil {
			m.Logger.Warnf("Error while fetching user %+v. err=%s", userDetails.User, userDetails.Error.Message)
			continue
		}

		user := emailUserMap[userDetails.User.Mail]
		if user == nil {
			m.Logger.Warnf("Couldn't find user with email %s.", userDetails.User.Mail)
			continue
		}

		err := m.CompleteUserConnect(user.Id, user.Timezone, userDetails.User)
		if err != nil {
			m.Logger.Warnf("Error connecting user with email %s. err=%s", user.Email, err.Error())
			continue
		}
	}

	return nil
}

func (m *mscalendar) CompleteOAuth2(authedUserID string) error {
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

	client := m.Remote.MakeClient(ctx)
	user, userErr := m.PluginAPI.GetMattermostUser(authedUserID)
	if userErr != nil {
		return userErr
	}

	me, err := client.GetMe(user.Email)
	if err != nil {
		return err
	}

	return m.CompleteUserConnect(authedUserID, user.Timezone, me)
}

func (m *mscalendar) CompleteUserConnect(authedUserID string, timezone model.StringMap, me *remote.User) error {
	_, err := m.Store.LoadMattermostUserID(me.ID)
	if err == nil {
		// Couldn't fetch connected MM account. Reject connect attempt.
		_, _ = m.Poster.DM(authedUserID, RemoteUserAlreadyConnectedNotFound, config.ApplicationName, me.Mail)
		return fmt.Errorf(RemoteUserAlreadyConnectedNotFound, config.ApplicationName, me.Mail)
	}

	u := &store.User{
		PluginVersion:    m.Config.PluginVersion,
		MattermostUserID: authedUserID,
		Remote:           me,
	}

	u.Settings.DailySummary = &store.DailySummaryUserSettings{
		PostTime: "8:00AM",
		Timezone: model.GetPreferredTimezone(timezone),
		Enable:   false,
	}

	err = m.Store.StoreUser(u)
	if err != nil {
		return err
	}

	err = m.Store.StoreUserInIndex(u)
	if err != nil {
		return err
	}

	_ = m.Welcomer.AfterSuccessfullyConnect(authedUserID, me.Mail)

	return nil
}
