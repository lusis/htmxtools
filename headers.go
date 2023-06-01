package htmxtools

// HTMXRequestHeader is a string type
type HTMXRequestHeader string

// String strings
func (rh HTMXRequestHeader) String() string {
	return string(rh)
}

// HTMXResponseHeader is a string type
type HTMXResponseHeader string

// String strings
func (rh HTMXResponseHeader) String() string {
	return string(rh)
}

// https://htmx.org/reference/#headers
// http headers are supposed to be case sensitive but
// in case something isn't behaving somewhere, we'll use the
// case from the project's page
const (
	hxHeaderPrefix = "HX-"
	// Request headers
	// HXRequestHeader is the header that signals an htmx request. Always true if called from htmx
	HXRequestHeader HTMXRequestHeader = hxHeaderPrefix + "Request"
	// BoostedRequest indicates that the request is via an element using hx-boost
	BoostedRequest HTMXRequestHeader = hxHeaderPrefix + "Boosted"
	// CurrentURLRequest the current URL of the browser
	CurrentURLRequest HTMXRequestHeader = hxHeaderPrefix + "Current-URL"
	// HistoryRestoreRequest is true if the request is for history restoration after a miss in the local history cache
	HistoryRestoreRequest HTMXRequestHeader = hxHeaderPrefix + "History-Restore-Request"
	// PromptRequest is the user response to an hx-prompt
	// https://htmx.org/attributes/hx-prompt/
	PromptRequest HTMXRequestHeader = hxHeaderPrefix + "Prompt"
	// TriggerRequest is the id of the target element if it exists
	TriggerRequest HTMXRequestHeader = hxHeaderPrefix + "Trigger"
	// TriggerNameRequest is the name of the triggered element if it exists
	TriggerNameRequest HTMXRequestHeader = hxHeaderPrefix + "Trigger-Name"
	// TargetRequest is the id of the target element if it exists
	TargetRequest HTMXRequestHeader = hxHeaderPrefix + "Target"

	// Response headers
	// LocationResponse Allows you to do a client-side redirect that does not do a full page reload
	// https://htmx.org/headers/hx-location/
	LocationResponse HTMXResponseHeader = hxHeaderPrefix + "Location"
	// PushURLResponse pushes a new url into the history stack
	// https://htmx.org/headers/hx-push-url/
	PushURLResponse HTMXResponseHeader = hxHeaderPrefix + "Push-Url"
	// RedirectResponse can be used to do a client-side redirect to a new location
	RedirectResponse HTMXResponseHeader = hxHeaderPrefix + "Redirect"
	// RefreshResponse if set to “true” the client side will do a a full refresh of the page
	RefreshResponse HTMXResponseHeader = hxHeaderPrefix + "Refresh"
	// ReplaceURLResponse replaces the current URL in the location bar
	// https://htmx.org/headers/hx-replace-url/
	ReplaceURLResponse HTMXResponseHeader = hxHeaderPrefix + "Replace-Url"
	// ReswapResponse Allows you to specify how the response will be swapped. See hx-swap for possible values
	ReswapResponse HTMXResponseHeader = hxHeaderPrefix + "Reswap"
	// RetargetResponse A CSS selector that updates the target of the content update to a different element on the page
	RetargetResponse HTMXResponseHeader = hxHeaderPrefix + "Retarget"
	// TriggerResponse allows you to trigger client side events, see the documentation for more info
	// https://htmx.org/headers/hx-trigger/
	TriggerResponse HTMXResponseHeader = hxHeaderPrefix + "Trigger"
	// TriggerAfterSettleResponse allows you to trigger client side events, see the documentation for more info
	// https://htmx.org/headers/hx-trigger/
	TriggerAfterSettleResponse HTMXResponseHeader = hxHeaderPrefix + "Trigger-After-Settle"
	// TriggerAfterSwapResponse allows you to trigger client side events, see the documentation for more info
	// https://htmx.org/headers/hx-trigger/
	TriggerAfterSwapResponse HTMXResponseHeader = hxHeaderPrefix + "Trigger-After-Swap"
)
