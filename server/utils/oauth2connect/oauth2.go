// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package oauth2connect

type App interface {
	InitOAuth2(mattermostUserID string) (string, error)
	CompleteOAuth2(mattermostUserID string) error
}

// type oa struct {
// 	app App
// }

// func Init(h *httputils.Handler, app App) {
// 	oa := &oa{
// 		app: app,
// 	}

// 	oauth2Router := h.Router.PathPrefix("/oauth2").Subrouter()
// 	oauth2Router.HandleFunc("/connect", oa.oauth2Connect).Methods("GET")
// 	oauth2Router.HandleFunc("/complete", oa.oauth2Complete).Methods("GET")
// }
