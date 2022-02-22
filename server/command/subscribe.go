// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) subscribe(parameters ...string) (string, bool, error) {
	_, err := c.MSCalendar.LoadMyEventSubscription()
	if err == nil {
		return "You are already subscribed to events.", false, nil
	}

	_, err = c.MSCalendar.CreateMyEventSubscription()
	if err != nil {
		return "", false, err
	}
	return "You are now subscribed to events.", false, nil
}
