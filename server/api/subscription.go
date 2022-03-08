package api

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/mscalendar"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
)

func (api *api) syncActionSubscription(w http.ResponseWriter, req *http.Request) {
	go func() {
		err := mscalendar.New(api.Env, "").SyncUserSubscriptions()
		if err != nil && err != store.ErrNotFound {
			api.Logger.Errorf("Error: Failed to sync user subscriptions: %s", err.Error())
			return
		}
	}()
}
