// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package oauth2connect

import (
	"net/http"
)

func (oa *oa) oauth2Connect(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get("Mattermost-User-ID")
	if mattermostUserID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// redirectURL, err := oa.app.InitOAuth2(mattermostUserID)
	// if err != nil {
	// 	httputils.WriteInternalServerError(w, err)
	// }
	// http.Redirect(w, r, redirectURL, http.StatusFound)

	oa.oauth2Complete(w, r)
}
