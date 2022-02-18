// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/httputils"
)

func (api *api) notification(w http.ResponseWriter, req *http.Request) {
	notification := api.Env.Remote.HandleWebhook(w, req)
	err := api.NotificationProcessor.Enqueue(&remote.Notification{
		ChangeType:     notification.ChangeType,
		SubscriptionID: notification.SubscriptionID,
		EventID:        notification.EventID,
	})
	if err != nil {
		httputils.WriteInternalServerError(w, err)
		return
	}
}
