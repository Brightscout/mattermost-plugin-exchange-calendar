// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"encoding/json"

	"io/ioutil"
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

type webhook struct {
	ChangeType     string `json:"changeType"`
	SubscriptionID string `json:"subscriptionId"`
	EventID        string `json:"eventId"`
}

func (r *impl) HandleWebhook(w http.ResponseWriter, req *http.Request) *remote.Notification {
	rawData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.logger.Infof("ews: failed to process webhook: `%v`.", err)
		return nil
	}

	var webhookResponse webhook
	err = json.Unmarshal(rawData, &webhookResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		r.logger.Infof("ews: failed to process webhook: `%v`.", err)
		return nil
	}

	n := &remote.Notification{
		ChangeType:     webhookResponse.ChangeType,
		SubscriptionID: webhookResponse.SubscriptionID,
		EventID:        webhookResponse.EventID,
	}

	w.WriteHeader(http.StatusAccepted)
	return n
}
