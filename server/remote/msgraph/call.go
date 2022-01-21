// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/pkg/errors"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func (c *client) CallJSON(method, path string, in, out interface{}) (responseData []byte, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, err
	}
	return c.call(method, path, contentType, buf, out)
}

func (c *client) CallFormPost(method, path string, in url.Values, out interface{}) (responseData []byte, err error) {
	contentType := "application/x-www-form-urlencoded"
	buf := strings.NewReader(in.Encode())
	return c.call(method, path, contentType, buf, out)
}

func (c *client) call(method, path, contentType string, inBody io.Reader, out interface{}) (responseData []byte, err error) {
	errContext := fmt.Sprintf("msgraph: Call failed: method:%s, path:%s", method, path)
	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, errors.WithMessage(err, errContext)
	}

	if pathURL.Scheme == "" || pathURL.Host == "" {
		var baseURL *url.URL
		baseURL, err = url.Parse(c.conf.ExchangeServerBaseURL)
		if err != nil {
			return nil, errors.WithMessage(err, errContext)
		}
		if path[0] != '/' {
			path = "/" + path
		}
		path = baseURL.String() + path
	}

	req, err := http.NewRequest(method, path, inBody)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if c.ctx != nil {
		req = req.WithContext(c.ctx)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if out != nil {
			err = json.Unmarshal(responseData, out)
			if err != nil {
				return responseData, err
			}
		}
		return responseData, nil

	case http.StatusNoContent:
		return nil, nil
	}

	errResp := msgraph.ErrorResponse{Response: resp}
	err = json.Unmarshal(responseData, &errResp)
	if err != nil {
		return responseData, errors.WithMessagef(err, "status: %s", resp.Status)
	}
	if err != nil {
		return responseData, err
	}
	return responseData, &errResp
}

func (c *client) GetEndpointURL(email, path string) (string, error) {
	endpointURL, err := url.Parse(strings.TrimSpace(fmt.Sprintf("%s%s", c.conf.ExchangeServerBaseURL, path)))
	if err != nil {
		return "", err
	}
	params := url.Values{}
	if email != "" {
		params.Add(config.EmailKey, email)
	}
	endpointURL.RawQuery = params.Encode()

	return endpointURL.String(), nil
}
