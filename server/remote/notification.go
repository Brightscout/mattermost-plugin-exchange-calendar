// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

type Notification struct {
	// Notification type
	ChangeType string

	// The (remote) subscription ID the notification is for
	SubscriptionID string

	// Set if subscription renewal is recommended. The date/time logic is
	// internal to the remote implementation. The handler is to call
	// RenewSubscription() as applicable, with the appropriate user credentials.
	RecommendRenew bool

	EventID string

	// Notification data
	Subscription        *Subscription
	SubscriptionCreator *User
	Event               *Event
}
