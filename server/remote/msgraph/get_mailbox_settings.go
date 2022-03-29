// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

func (c *client) GetMailboxSettings(remoteUserID string) (*remote.MailboxSettings, error) {
	return &remote.MailboxSettings{}, nil
}
