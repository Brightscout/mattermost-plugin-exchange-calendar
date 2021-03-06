// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "net/http"

func (api *api) getAuthorized(w http.ResponseWriter, r *http.Request) {
	// if we've made it here, we're authorized.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"is_authorized": true}`))
}
