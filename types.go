package htmxtools

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type htmxContextKey string

var requestContextKey = htmxContextKey("htmxRequest")

// HTMXRequest represents the htmx elements of an [http.Request]
// fields may be empty strings
type HTMXRequest struct {
	// https://htmx.org/attributes/hx-boost/
	Boosted        bool
	CurrentURL     string
	HistoryRestore bool
	// https://htmx.org/attributes/hx-prompt/
	Prompt      string
	Target      string
	TriggerName string
	Trigger     string
}

// ToContext adds the HTMXRequest details to the provided parent context
func (hr *HTMXRequest) ToContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	return context.WithValue(ctx, requestContextKey, hr)
}

// RequestFromContext parses the htmx request from the provided context
func RequestFromContext(ctx context.Context) *HTMXRequest {
	res, ok := ctx.Value(requestContextKey).(*HTMXRequest)
	if !ok {
		return nil
	}
	return res
}

// HTMXResponse represents the htmx elements of an [http.Response]
type HTMXResponse struct {
	// https://htmx.org/headers/hx-location/
	Location string
	// https://htmx.org/headers/hx-push-url/
	PushURL  string
	Redirect string
	Refresh  bool
	// https://htmx.org/headers/hx-replace-url/
	ReplaceURL string
	// https://htmx.org/attributes/hx-swap/
	Reswap HXSwap
	// https://htmx.org/headers/hx-trigger/
	Trigger string
	// https://htmx.org/headers/hx-trigger/
	TriggerAfterSettle string
	// https://htmx.org/headers/hx-trigger/
	TriggerAfterSwap string
}

// AddHeaders adds the provided htmx headers to the request
func (hr *HTMXResponse) AddHeaders(r *http.Request, headers ...map[HTMXResponseHeader]string) error {
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
