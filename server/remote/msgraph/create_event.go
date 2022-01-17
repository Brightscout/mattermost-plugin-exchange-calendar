// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/pkg/errors"
)

// CreateEvent creates a calendar event
func (c *client) CreateEvent(remoteUserEmail string, in *remote.Event) (*remote.Event, error) {
	var out = &remote.Event{}
	url, err := c.GetEndpointURL(remoteUserEmail, config.PathEvent)
	if err != nil {
		return nil, errors.Wrap(err, "ews CreateEvent")
	}
	_, err = c.CallJSON(http.MethodPost, url, in, out)
	if err != nil {
		return nil, errors.Wrap(err, "ews CreateEvent")
	}
	return out, nil
}
