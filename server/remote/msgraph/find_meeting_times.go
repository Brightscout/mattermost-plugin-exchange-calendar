// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

// FindMeetingTimes finds meeting time suggestions for a calendar event
func (c *client) FindMeetingTimes(remoteUserID string, params *remote.FindMeetingTimesParameters) (*remote.MeetingTimeSuggestionResults, error) {
	meetingsOut := &remote.MeetingTimeSuggestionResults{}
	// TODO: Add FindMeetingTimes API
	// req := c.rbuilder.Users().ID(remoteUserID).FindMeetingTimes(nil).Request()
	// err := req.JSONRequest(c.ctx, http.MethodPost, "", &params, &meetingsOut)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "msgraph FindMeetingTimes")
	// }
	return meetingsOut, nil
}
