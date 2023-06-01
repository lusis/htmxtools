package htmxtools

import (
	"encoding/json"
	"net/http"
	"strings"
)

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

// AddToResponse adds the current state of the HTMXResponse to the http response headers
func (hr *HTMXResponse) AddToResponse(w http.ResponseWriter) error {
	if strings.TrimSpace(hr.Location) != "" {
		w.Header().Set(LocationResponse.String(), hr.Location)
	}
	if strings.TrimSpace(hr.PushURL) != "" {
		w.Header().Set(PushURLResponse.String(), hr.PushURL)
	}
	if strings.TrimSpace(hr.Redirect) != "" {
		w.Header().Set(RedirectResponse.String(), hr.Redirect)
	}
	if hr.Refresh {
		w.Header().Set(RefreshResponse.String(), "true")
	}
	if strings.TrimSpace(hr.ReplaceURL) != "" {
		w.Header().Set(ReplaceURLResponse.String(), hr.ReplaceURL)
	}
	if strings.TrimSpace(hr.Trigger) != "" {
		w.Header().Set(TriggerResponse.String(), hr.Trigger)
	}
	if strings.TrimSpace(hr.TriggerAfterSettle) != "" {
		w.Header().Set(TriggerAfterSettleResponse.String(), hr.TriggerAfterSettle)
	}
	if hr.Reswap != 0 {
		w.Header().Set(ReswapResponse.String(), hr.Reswap.String())
	}
	if strings.TrimSpace(hr.TriggerAfterSwap) != "" {
		w.Header().Set(TriggerAfterSwapResponse.String(), hr.TriggerAfterSwap)
	}
	return nil
}

// HXLocationResponse represents the structured format of an hx-location header described here:
// https://htmx.org/headers/hx-location/
type HXLocationResponse struct {
	Path    string                 `json:"path"`
	Source  string                 `json:"source,omitempty"`
	Event   string                 `json:"event,omitempty"`
	Handler string                 `json:"handler,omitempty"`
	Target  string                 `json:"target,omitempty"`
	Swap    HXSwap                 `json:"swap,omitempty"`
	Values  map[string]interface{} `json:"value,omitempty"`
	Headers map[string]string      `json:"headers,omitempty"`
}

// String strings
func (hxl HXLocationResponse) String() string {
	j, _ := json.Marshal(hxl)
	return string(j)
}
