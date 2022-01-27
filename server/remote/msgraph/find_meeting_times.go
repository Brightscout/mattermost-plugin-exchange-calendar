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

// FindMeetingTimes finds meeting time suggestions for a calendar event
func (c *client) FindMeetingTimes(remoteUserEmail string, params *remote.FindMeetingTimesParameters) (*remote.MeetingTimeSuggestionResults, error) {
	meetingsOut := &remote.MeetingTimeSuggestionResults{}
	url, err := c.GetEndpointURL(remoteUserEmail, fmt.Sprintf("%s%s", config.PathCalendar, config.PathFindMeetingTimes))
	if err != nil {
		return nil, errors.Wrap(err, "ews FindMeetingTimes")
	}
	_, err = c.CallJSON(http.MethodPost, url, params, &meetingsOut)
	if err != nil {
		return nil, errors.Wrap(err, "ews FindMeetingTimes")
	}
	return meetingsOut, nil
}
