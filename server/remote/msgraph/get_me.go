// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

func (c *client) GetMe(remoteUserEmail string) (*remote.User, error) {
	var remoteUser remote.User
	path, err := c.GetEndpointURL(config.PathMe, &remoteUserEmail)
	if err != nil {
		return nil, errors.Wrap(err, "ews GetMe")
	}

	_, err = c.CallJSON(http.MethodGet, path, nil, &remoteUser)
	if err != nil {
		return nil, errors.Wrap(err, "ews GetMe")
	}

	if remoteUser.Mail == "" {
		return nil, errors.New("user has no email address. Make sure the Microsoft account is associated to an Outlook product")
	}
	if remoteUser.DisplayName == "" {
		return nil, errors.New("user has no Display Name")
	}

	return &remoteUser, nil
}
