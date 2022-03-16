package utils

import (
	"encoding/json"
	"io"

	"github.com/mattermost/mattermost-server/v6/model"
)

func PostActionIntegrationRequestFromJson(data io.Reader) *model.PostActionIntegrationRequest {
	var o *model.PostActionIntegrationRequest
	err := json.NewDecoder(data).Decode(&o)
	if err != nil {
		return nil
	}
	return o
}

func ResponseToJson(response model.PostActionIntegrationResponse) []byte {
	byteResponse, _ := json.Marshal(response)
	return byteResponse
}
