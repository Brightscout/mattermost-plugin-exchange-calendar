// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"encoding/xml"

	"io/ioutil"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

type WebhookResponseEnvelope struct {
	XMLName xml.Name
	Body    Body
}

type Body struct {
	XMLName          xml.Name
	SendNotification SendNotification `xml:"SendNotification"`
}

type SendNotification struct {
	XMLName          xml.Name
	ResponseMessages ResponseMessages `xml:"ResponseMessages"`
}

type ResponseMessages struct {
	XMLName                         xml.Name
	SendNotificationResponseMessage SendNotificationResponseMessage `xml:"SendNotificationResponseMessage"`
}

type SendNotificationResponseMessage struct {
	XMLName       xml.Name
	ResponseClass string       `xml:"ResponseClass,attr"`
	ResponseCode  string       `xml:"ResponseCode"`
	Notification  Notification `xml:"Notification"`
}

type Notification struct {
	XMLName        xml.Name
	SubscriptionID string       `xml:"SubscriptionId"`
	CreatedEvent   CreatedEvent `xml:"CreatedEvent"`
	StatusEvent    *StatusEvent `xml:"StatusEvent"`
}

type CreatedEvent struct {
	XMLName xml.Name
	Item    Item `xml:"ItemId"`
}

type StatusEvent struct {
	XMLName xml.Name
}

type Item struct {
	XMLName xml.Name
	EventID string `xml:"Id,attr"`
}

func (r *impl) HandleWebhook(w http.ResponseWriter, req *http.Request) (bool, *remote.Notification, error) {
	rawData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.logger.Infof("ews: failed to process webhook: `%v`.", err)
		return false, nil, err
	}

	var webhookResponse WebhookResponseEnvelope
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
	statusEvent := webhookResponse.Body.SendNotification.ResponseMessages.SendNotificationResponseMessage.Notification.StatusEvent

	return statusEvent != nil, n, nil
}
