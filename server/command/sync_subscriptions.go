package command

func (c *Command) syncSubscriptions(parameters ...string) (string, bool, error) {

	// Sync the subscription asynchronously as it will take time
	go func() {
		_ = c.MSCalendar.SyncUserSubscriptions()
	}()

	return "Syncing subscriptions started successfully.", false, nil
}
