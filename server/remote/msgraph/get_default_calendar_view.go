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
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
)

type calendarViewSingleRequest struct {
	ID            string    `json:"id"`
	StartDateTime time.Time `json:"startDateTime"`
	EndDateTime   time.Time `json:"endDateTime"`
}

type calendarViewBatchRequest struct {
	Requests []*calendarViewSingleRequest `json:"requests"`
}

type calendarViewSingleResponse struct {
	ID     string           `json:"id"`
	Events []*remote.Event  `json:"events"`
	Error  *remote.APIError `json:"error,omitempty"`
}

type calendarViewBatchResponse struct {
	Responses []*calendarViewSingleResponse `json:"responses"`
}

func (c *client) GetDefaultCalendarView(remoteUserEmail string, start, end time.Time) ([]*remote.Event, error) {
	var out []*remote.Event
	url, err := c.GetEndpointURL(config.PathEvent, &remoteUserEmail)
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
	requests := []*calendarViewSingleRequest{}
	for _, params := range allParams {
		req := &calendarViewSingleRequest{
			ID:            params.RemoteUserID,
			StartDateTime: params.StartTime,
			EndDateTime:   params.EndTime,
		}
		requests = append(requests, req)
	}

	batchRequests := prepareEventBatchRequests(requests)
	var batchResponses []*calendarViewBatchResponse
	for _, req := range batchRequests {
		batchRes := &calendarViewBatchResponse{}
		err := c.GetEventsBatchRequest(req, batchRes)
		if err != nil {
			return nil, errors.Wrap(err, "ews ViewCalendar batch request")
		}

		batchResponses = append(batchResponses, batchRes)
	}

	result := []*remote.ViewCalendarResponse{}
	for _, batchRes := range batchResponses {
		for _, res := range batchRes.Responses {
			viewCalRes := &remote.ViewCalendarResponse{
				RemoteUserID: res.ID,
				Events:       res.Events,
				Error:        res.Error,
			}
			result = append(result, viewCalRes)
		}
	}

	return result, nil
}

func prepareEventBatchRequests(requests []*calendarViewSingleRequest) []calendarViewBatchRequest {
	numOfBatches := utils.GetTotalNumberOfBatches(len(requests), maxNumRequestsPerBatch)
	result := []calendarViewBatchRequest{}

	for i := 0; i < numOfBatches; i++ {
		startIdx := i * maxNumRequestsPerBatch
		endIdx := startIdx + maxNumRequestsPerBatch
		// In case of last batch endIdx will be equal to length of requests
		if i == numOfBatches-1 {
			endIdx = len(requests)
		}

		slice := requests[startIdx:endIdx]
		batchReq := calendarViewBatchRequest{Requests: slice}
		result = append(result, batchReq)
	}

	return result
}

func (c *client) GetEventsBatchRequest(req calendarViewBatchRequest, out interface{}) error {
	url, err := c.GetEndpointURL(config.PathBatchEvent, nil)
	if err != nil {
		return errors.Wrap(err, "ews GetEventsBatchRequest")
	}
	_, err = c.CallJSON(http.MethodPost, url, req, out)
	if err != nil {
		return errors.Wrap(err, "ews GetEventsBatchRequest")
	}

	return nil
}

func getQueryParamStringForCalendarView(start, end time.Time) string {
	q := url.Values{}
	q.Add("startDateTime", start.Format(time.RFC3339))
	q.Add("endDateTime", end.Format(time.RFC3339))
	return q.Encode()
}
