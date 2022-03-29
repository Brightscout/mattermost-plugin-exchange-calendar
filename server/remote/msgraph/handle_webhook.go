// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"encoding/xml"

	"io/ioutil"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

func (r *impl) HandleWebhook(w http.ResponseWriter, req *http.Request) (bool, *remote.Notification, error) {
	rawData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.logger.Infof("ews: failed to process webhook: `%v`.", err)
		return false, nil, err
	}

	var webhookResponse remote.WebhookResponseEnvelope
	err = xml.Unmarshal(rawData, &webhookResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		r.logger.Infof("ews: failed to process webhook: `%v`.", err)
		return false, nil, err
	}

	n := &remote.Notification{
		SubscriptionID: webhookResponse.Body.SendNotification.ResponseMessages.SendNotificationResponseMessage.Notification.SubscriptionID,
		EventID:        webhookResponse.Body.SendNotification.ResponseMessages.SendNotificationResponseMessage.Notification.CreatedEvent.Item.EventID,
	}

	// statusEvent indicates whether the webhook request from the Exchange server is for a status check or it contains notification data
	statusEvent := webhookResponse.Body.SendNotification.ResponseMessages.SendNotificationResponseMessage.Notification.StatusEvent

	return statusEvent != nil, n, nil
}
