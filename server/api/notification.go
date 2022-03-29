// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"encoding/xml"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
)

func (api *api) notification(w http.ResponseWriter, req *http.Request) {
	isStatusCheck, notification, err := api.Env.Remote.HandleWebhook(w, req)
	if err != nil {
		return
	}

	subscriptionStatusText := config.SubscriptionStatusOK
	// Load subscription from store
	_, err = api.Store.LoadSubscription(notification.SubscriptionID)
	if err != nil {
		if err != store.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// If subscription doesn't exist in store, unsubscribe from notifications
		subscriptionStatusText = config.SubscriptionStatusUnsubscribe
		isStatusCheck = true
	}

	if !isStatusCheck {
		err = api.NotificationProcessor.Enqueue(&remote.Notification{
			ChangeType:     notification.ChangeType,
			SubscriptionID: notification.SubscriptionID,
			EventID:        notification.EventID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	soapResponse := remote.SOAPEnvelope{
		Body: remote.SOAPBody{
			Content: remote.SendNotificationResult{
				SubscriptionStatus: subscriptionStatusText,
			},
		},
	}
	marshalledResponse, err := xml.MarshalIndent(soapResponse, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(marshalledResponse)
}
