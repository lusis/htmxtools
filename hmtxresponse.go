package htmxtools

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
