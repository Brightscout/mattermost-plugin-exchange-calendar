// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

func (c *client) GetMailboxSettings(remoteUserID string) (*remote.MailboxSettings, error) {
	// TODO: Add GetMailboxSettings API
	// u := c.rbuilder.Users().ID(remoteUserID).URL() + "/mailboxSettings"
	out := &remote.MailboxSettings{}

	// _, err := c.CallJSON(http.MethodGet, u, nil, out)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "msgraph GetMailboxSettings")
	// }
	return out, nil
}
