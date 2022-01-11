// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
)

type calendarViewResponse struct {
	Value []*remote.Event  `json:"value,omitempty"`
	Error *remote.APIError `json:"error,omitempty"`
}

type calendarViewSingleResponse struct {
	ID      string               `json:"id"`
	Status  int                  `json:"status"`
	Body    calendarViewResponse `json:"body"`
	Headers map[string]string    `json:"headers"`
}

type calendarViewBatchResponse struct {
	Responses []*calendarViewSingleResponse `json:"responses"`
}

func (c *client) GetDefaultCalendarView(remoteUserEmail string, start, end time.Time) ([]*remote.Event, error) {
	var out []*remote.Event
	url, err := c.GetEndpointURL(remoteUserEmail, config.PathEvent)
	if err != nil {
		return nil, errors.Wrap(err, "ews GetDefaultCalendarView")
	}
	url = fmt.Sprintf("%s&%s", url, getQueryParamStringForCalendarView(start, end))
	_, err = c.CallJSON(http.MethodGet, url, nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "ews GetDefaultCalendarView")
	}

	return out, nil
}

func (c *client) DoBatchViewCalendarRequests(allParams []*remote.ViewCalendarParams) ([]*remote.ViewCalendarResponse, error) {
	requests := []*singleRequest{}
	for _, params := range allParams {
		u := getCalendarViewURL(params)
		req := &singleRequest{
			ID:      params.RemoteUserID,
			URL:     u,
			Method:  http.MethodGet,
			Headers: map[string]string{},
		}
		requests = append(requests, req)
	}

	batchRequests := prepareBatchRequests(requests)
	var batchResponses []*calendarViewBatchResponse
	for _, req := range batchRequests {
		batchRes := &calendarViewBatchResponse{}
		err := c.batchRequest(req, batchRes)
		if err != nil {
			return nil, errors.Wrap(err, "msgraph ViewCalendar batch request")
		}

		batchResponses = append(batchResponses, batchRes)
	}

	result := []*remote.ViewCalendarResponse{}
	for _, batchRes := range batchResponses {
		for _, res := range batchRes.Responses {
			viewCalRes := &remote.ViewCalendarResponse{
				RemoteUserID: res.ID,
				Events:       res.Body.Value,
				Error:        res.Body.Error,
			}
			result = append(result, viewCalRes)
		}
	}

	return result, nil
}

func getCalendarViewURL(params *remote.ViewCalendarParams) string {
	paramStr := "?" + getQueryParamStringForCalendarView(params.StartTime, params.EndTime)
	return "/Users/" + params.RemoteUserID + "/calendarView" + paramStr
}

func getQueryParamStringForCalendarView(start, end time.Time) string {
	q := url.Values{}
	q.Add("startDateTime", start.Format(time.RFC3339))
	q.Add("endDateTime", end.Format(time.RFC3339))
	return q.Encode()
}
