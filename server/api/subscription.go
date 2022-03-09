package api

import (
	"net/http"
	"strconv"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/mscalendar"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/httputils"
	"github.com/gorilla/mux"
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

func (api *api) getSubscritpionByID(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	subscriptionID := params["subscriptionID"]
	decodedID, err := utils.DecodeString(subscriptionID)
	if err != nil {
		httputils.WriteInternalServerError(w, err)
		return
	}
	_, err = mscalendar.New(api.Env, "").GetSubscritpionByID(decodedID)
	isSubscribed := true
	if err != nil {
		if err != store.ErrNotFound {
			httputils.WriteInternalServerError(w, err)
			return
		}
		isSubscribed = false
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(strconv.FormatBool(isSubscribed)))
}
