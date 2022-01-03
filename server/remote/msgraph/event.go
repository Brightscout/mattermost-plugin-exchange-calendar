// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

func (c *client) GetEvent(remoteUserID, eventID string) (*remote.Event, error) {
	e := &remote.Event{}
	// TODO: Add GetEvent API
	// err := c.rbuilder.Users().ID(remoteUserID).Events().ID(eventID).Request().JSONRequest(
	// 	c.ctx, http.MethodGet, "", nil, &e)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "msgraph GetEvent")
	// }
	return e, nil
}

func (c *client) AcceptEvent(remoteUserID, eventID string) error {
	// TODO: Add AcceptEvent API
	// dummy := &msgraph.EventAcceptRequestParameter{}
	// err := c.rbuilder.Users().ID(remoteUserID).Events().ID(eventID).Accept(dummy).Request().Post(c.ctx)
	// if err != nil {
	// 	return errors.Wrap(err, "msgraph Accept Event")
	// }
	return nil
}

func (c *client) DeclineEvent(remoteUserID, eventID string) error {
	// TODO: Add DeclineEvent API
	// dummy := &msgraph.EventDeclineRequestParameter{}
	// err := c.rbuilder.Users().ID(remoteUserID).Events().ID(eventID).Decline(dummy).Request().Post(c.ctx)
	// if err != nil {
	// 	return errors.Wrap(err, "msgraph DeclineEvent")
	// }
	return nil
}

func (c *client) TentativelyAcceptEvent(remoteUserID, eventID string) error {
	// TODO: Add TentativelyAcceptEvent API
	// dummy := &msgraph.EventTentativelyAcceptRequestParameter{}
	// err := c.rbuilder.Users().ID(remoteUserID).Events().ID(eventID).TentativelyAccept(dummy).Request().Post(c.ctx)
	// if err != nil {
	// 	return errors.Wrap(err, "msgraph TentativelyAcceptEvent")
	// }
	return nil
}
