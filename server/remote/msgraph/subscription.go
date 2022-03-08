// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
)

func (c *client) GetWebhookNotificationURL() string {
	return fmt.Sprintf("%s%s%s", c.conf.MattermostSiteURL, c.conf.PluginURLPath, config.FullPathEventNotification)
}

func (c *client) CreateMySubscription(remoteUserEmail string, notificationURL string) (*remote.Subscription, error) {
	sub := &remote.Subscription{
		WebhookNotificationUrl: c.GetWebhookNotificationURL(),
	}

	path, err := c.GetEndpointURL(fmt.Sprintf("%s%s", config.PathNotification, config.PathSubscribe), &remoteUserEmail)
	if err != nil {
		return nil, errors.Wrap(err, "ews Subscribe")
	}

	_, err = c.CallJSON(http.MethodPost, path, &sub, &sub)
	if err != nil {
		return nil, errors.Wrap(err, "ews Subscribe")
	}

	sub.CreatorID = remoteUserEmail

	c.Logger.With(bot.LogContext{
		"subscriptionID": sub.ID,
	}).Debugf("ews: created subscription.")

	return sub, nil
}

func (c *client) DeleteSubscription(remoteUserEmail, subscriptionID string) error {
	sub := &remote.Subscription{
		ID: subscriptionID,
	}

	path, err := c.GetEndpointURL(fmt.Sprintf("%s%s", config.PathNotification, config.PathUnsubscribe), &remoteUserEmail)
	if err != nil {
		return errors.Wrap(err, "ews DeleteSubscription")
	}
	_, err = c.CallJSON(http.MethodPost, path, &sub, nil)
	if err != nil {
		return errors.Wrap(err, "ews DeleteSubscription")
	}

	c.Logger.With(bot.LogContext{
		"subscriptionID": subscriptionID,
	}).Debugf("ews: deleted subscription.")

	return nil
}

func (c *client) DoBatchSubscriptionRequests(requests []remote.SubscriptionBatchSingleRequest) ([]*remote.SubscriptionBatchSingleResponse, error) {
	batchRequests := prepareSubscriptionBatchRequests(requests)
	var batchResponses []*remote.SubscriptionBatchSingleResponse
	for _, req := range batchRequests {
		batchResponse := []*remote.SubscriptionBatchSingleResponse{}
		err := c.GetSubscriptionsBatchRequest(req, &batchResponse)
		if err != nil {
			return nil, errors.Wrap(err, "ews Subscription batch request")
		}

		batchResponses = append(batchResponses, batchResponse...)
	}

	return batchResponses, nil
}

func prepareSubscriptionBatchRequests(requests []remote.SubscriptionBatchSingleRequest) [][]remote.SubscriptionBatchSingleRequest {
	numOfBatches := utils.GetTotalNumberOfBatches(len(requests), maxNumRequestsPerBatch)
	result := [][]remote.SubscriptionBatchSingleRequest{}
	for i := 0; i < numOfBatches; i++ {
		startIdx := i * maxNumRequestsPerBatch
		endIdx := startIdx + maxNumRequestsPerBatch
		// In case of last batch endIdx will be equal to length of requests
		if i == numOfBatches - 1 {
			endIdx = len(requests)
		}

		result = append(result, requests[startIdx:endIdx])
	}

	return result
}

func (c *client) GetSubscriptionsBatchRequest(req []remote.SubscriptionBatchSingleRequest, out interface{}) error {
	url, err := c.GetEndpointURL(config.PathBatchSubscription, nil)
	if err != nil {
		return errors.Wrap(err, "ews GetSubscriptionsBatchRequest")
	}
	_, err = c.CallJSON(http.MethodPost, url, req, out)
	if err != nil {
		return errors.Wrap(err, "ews GetSubscriptionsBatchRequest")
	}

	return nil
}
