package htmxtools

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
