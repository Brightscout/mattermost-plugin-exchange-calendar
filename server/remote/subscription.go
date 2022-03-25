// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

type Subscription struct {
	ID                     string `json:"subscriptionId"`
	WebhookNotificationUrl string `json:"webhookNotificationUrl,omitempty"`
	CreatorID              string `json:"creatorId,omitempty"`
}

type SubscriptionBatchSingleRequest struct {
	Email        string       `json:"email"`
	Subscription Subscription `json:"subscription"`
}

type SubscriptionBatchSingleResponse struct {
	Email        string         `json:"email"`
	Subscription *Subscription  `json:"subscription"`
	Error        *ErrorResponse `json:"error"`
}
