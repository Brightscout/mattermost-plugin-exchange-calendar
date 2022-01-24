// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"fmt"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/pkg/errors"
)

func (c *client) GetEvent(remoteUserEmail, eventID string) (*remote.Event, error) {
	e := &remote.Event{}
	url, err := c.GetEndpointURL(fmt.Sprintf("%s/%s", config.PathEvent, eventID), &remoteUserEmail)
	if err != nil {
		return nil, errors.Wrap(err, "ews GetEvent")
	}
	_, err = c.CallJSON(http.MethodGet, url, nil, &e)
	if err != nil {
		return nil, errors.Wrap(err, "ews GetEvent")
	}
	return e, nil
}

func (c *client) AcceptEvent(remoteUserEmail, eventID string) error {
	url, err := c.GetEndpointURL(fmt.Sprintf("%s%s/%s", config.PathEvent, config.PathAccept, eventID), &remoteUserEmail)
	if err != nil {
		return errors.Wrap(err, "ews AcceptEvent")
	}
	_, err = c.CallJSON(http.MethodGet, url, nil, nil)
	if err != nil {
		return errors.Wrap(err, "ews AcceptEvent")
	}
	return nil
}

func (c *client) DeclineEvent(remoteUserEmail, eventID string) error {
	url, err := c.GetEndpointURL(fmt.Sprintf("%s%s/%s", config.PathEvent, config.PathDecline, eventID), &remoteUserEmail)
	if err != nil {
		return errors.Wrap(err, "ews DeclineEvent")
	}
	_, err = c.CallJSON(http.MethodGet, url, nil, nil)
	if err != nil {
		return errors.Wrap(err, "ews DeclineEvent")
	}
	return nil
}

func (c *client) TentativelyAcceptEvent(remoteUserEmail, eventID string) error {
	url, err := c.GetEndpointURL(fmt.Sprintf("%s%s/%s", config.PathEvent, config.PathTentative, eventID), &remoteUserEmail)
	if err != nil {
		return errors.Wrap(err, "ews TentativelyAcceptEvent")
	}
	_, err = c.CallJSON(http.MethodGet, url, nil, nil)
	if err != nil {
		return errors.Wrap(err, "ews TentativelyAcceptEvent")
	}
	return nil
}
