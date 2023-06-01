package htmxtools

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRequest(t *testing.T) {
	t.Parallel()
	type testcase struct {
		httpRequest *http.Request
		htmxRequest *HTMXRequest
	}
	validRequest, _ := http.NewRequest(http.MethodGet, "localhost", nil)
	validRequest.Header.Set(HXRequestHeader.String(), "true")
	fullRequest := validRequest.Clone(context.TODO())
	fullRequest.Header.Add(BoostedRequest.String(), "true")
	fullRequest.Header.Add(CurrentURLRequest.String(), "localhost")
	fullRequest.Header.Add(HistoryRestoreRequest.String(), "true")
	fullRequest.Header.Add(PromptRequest.String(), "did you do it?")
	fullRequest.Header.Add(TriggerRequest.String(), "add-thing")
	fullRequest.Header.Add(TargetRequest.String(), "target-div")
	fullRequest.Header.Add(TriggerNameRequest.String(), "thing-id")
	testcases := map[string]testcase{
		"nil":          {httpRequest: &http.Request{}},
		"validrequest": {httpRequest: validRequest, htmxRequest: &HTMXRequest{}},
		"fullrequest": {httpRequest: fullRequest, htmxRequest: &HTMXRequest{
			Boosted:        true,
			CurrentURL:     "localhost",
			HistoryRestore: true,
			Prompt:         "did you do it?",
			Target:         "target-div",
			Trigger:        "add-thing",
			TriggerName:    "thing-id",
		}},
	}
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res := ParseRequest(tc.httpRequest)
			if tc.htmxRequest == nil {
				require.Nil(t, res)
			} else {
				require.NotNil(t, res)
				require.Equal(t, tc.htmxRequest.Boosted, res.Boosted, "boosted should match")
				require.Equal(t, tc.htmxRequest.CurrentURL, res.CurrentURL, "current url should match")
				require.Equal(t, tc.htmxRequest.HistoryRestore, res.HistoryRestore, "history restore should match")
				require.Equal(t, tc.htmxRequest.Prompt, res.Prompt, "prompt should match")
				require.Equal(t, tc.htmxRequest.Target, res.Target, "target should match")
				require.Equal(t, tc.htmxRequest.Trigger, res.Trigger, "trigger should match")
				require.Equal(t, tc.htmxRequest.TriggerName, res.TriggerName, "trigger name should match")
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	t.Parallel()
	htmxReq := &HTMXRequest{
		Boosted:        true,
		CurrentURL:     "localhost",
		HistoryRestore: true,
		Prompt:         "did you do it?",
		Target:         "target-div",
		Trigger:        "add-thing",
		TriggerName:    "thing-id",
	}
	var extractedRequest *HTMXRequest
	// need a func to wrap with the middleware and capture the context details
	testHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		extractedRequest = RequestFromContext(r.Context())
	}
	server := httptest.NewServer(Middleware(http.HandlerFunc(testHandlerFunc)))
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err, "new request should not error")
	require.NotNil(t, req, "http request should not be nil")
	// do the request without the headers first
	nres, nerr := http.DefaultClient.Do(req)
	require.NoError(t, nerr, "request should not error")
	require.NotNil(t, nres, "result should not be nil")
	require.Nil(t, extractedRequest, "extracted result should be nil")
	// add the htmx header
	req.Header.Set(HXRequestHeader.String(), "true")
	req.Header.Add(BoostedRequest.String(), "true")
	req.Header.Add(CurrentURLRequest.String(), "localhost")
	req.Header.Add(HistoryRestoreRequest.String(), "true")
	req.Header.Add(PromptRequest.String(), "did you do it?")
	req.Header.Add(TriggerRequest.String(), "add-thing")
	req.Header.Add(TargetRequest.String(), "target-div")
	req.Header.Add(TriggerNameRequest.String(), "thing-id")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "making request should not error")
	require.NotNil(t, res, "http request result should not be nil")
	require.NotNil(t, extractedRequest, "request from context should not be nil")
	require.Equal(t, htmxReq.Boosted, extractedRequest.Boosted)
}

// addHeaders adds hmtx headers to the provided http request for testing
func (hr *HTMXResponse) addHeaders(r *http.Request, headers ...map[HTMXRequestHeader]string) error { // nolint: unused
	for _, h := range headers {
		for k, v := range h {
			if strings.TrimSpace(r.Header.Get(k.String())) != "" {
				return fmt.Errorf("header already set: %s: %s", k, v)
			}
			r.Header.Set(k.String(), v)
		}
	}
	return nil
}
