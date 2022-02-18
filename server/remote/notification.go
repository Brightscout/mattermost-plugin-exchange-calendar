// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

type Notification struct {
	// Notification type
	ChangeType string

	// The (remote) subscription ID the notification is for
	SubscriptionID string

	EventID string

	// Notification data
	Subscription        *Subscription
	SubscriptionCreator *User
	Event               *Event
}
