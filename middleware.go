package htmxtools

import (
	"net/http"
	"strings"
)

// Middleware parses an htmx request and injects any htmx metadata in the context
// drop indicates if non-htmx requests should be dropped or not
func Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if res := parseRequest(r); res != nil {
			ctx := res.ToContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	}
}

// ParseRequest parses an [http.Request] for any htmx request headers and returns an [HTMXRequest]
// fields will still have to be checked for empty string at call sites
func parseRequest(r *http.Request) *HTMXRequest {
	isHTMXRequest := strings.TrimSpace(r.Header.Get(HXRequestHeader.String())) == "true"
	if !isHTMXRequest {
		return nil
	}
	res := &HTMXRequest{
		Boosted:        strings.TrimSpace(r.Header.Get(BoostedRequest.String())) == "true",
		CurrentURL:     strings.TrimSpace(r.Header.Get(CurrentURLRequest.String())),
		HistoryRestore: strings.TrimSpace(r.Header.Get(HistoryRestoreRequest.String())) == "true",
		Prompt:         strings.TrimSpace(r.Header.Get(PromptRequest.String())),
		Target:         strings.TrimSpace(r.Header.Get(TargetRequest.String())),
		TriggerName:    strings.TrimSpace(r.Header.Get(TriggerNameRequest.String())),
		Trigger:        strings.TrimSpace(r.Header.Get(TriggerRequest.String())),
	}

	return res
}
