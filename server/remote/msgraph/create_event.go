// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

// CreateEvent creates a calendar event
func (c *client) CreateEvent(remoteUserID string, in *remote.Event) (*remote.Event, error) {
	var out = remote.Event{}
	// TODO: Add CreateEvent API
	// err := c.rbuilder.Users().ID(remoteUserID).Events().Request().JSONRequest(c.ctx, http.MethodPost, "", &in, &out)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "msgraph CreateEvent")
	// }
	return &out, nil
}
