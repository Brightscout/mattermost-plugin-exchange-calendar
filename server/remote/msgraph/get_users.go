package msgraph

import (
	"net/http"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/pkg/errors"
)

func (c *client) GetUsers(emails []*string) ([]*remote.UserBatchSingleResponse, error) {
	url, err := c.GetEndpointURL(config.PathBatchUser, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ews GetUsers")
	}
	res := []*remote.UserBatchSingleResponse{}
	_, err = c.CallJSON(http.MethodPost, url, emails, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
