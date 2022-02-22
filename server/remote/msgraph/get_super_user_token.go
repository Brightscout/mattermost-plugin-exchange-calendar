// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"net/url"
)

type AuthResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func (c *client) GetSuperuserToken() (string, error) {
	params := map[string]string{
		"scope":      "https://graph.microsoft.com/.default",
		"grant_type": "client_credentials",
	}

	res := AuthResponse{}

	data := url.Values{}
	data.Set("client_id", params["client_id"])
	data.Set("scope", params["scope"])
	data.Set("client_secret", params["client_secret"])
	data.Set("grant_type", params["grant_type"])

	return res.AccessToken, nil
}
