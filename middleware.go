package htmxtools

import (
	"net/http"
	"strings"
)

// WrapFunc is middleware for inspecting http requests for htmx metadata
func WrapFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if res := ParseRequest(r); res != nil {
			ctx := res.ToContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Wrap is middleware for inspecting http requests for htmx metadata
func Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if res := ParseRequest(r); res != nil {
			ctx := res.ToContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ParseRequest parses an [http.Request] for any htmx request headers and returns an [HTMXRequest]
// fields will still have to be checked for empty string at call sites
func ParseRequest(r *http.Request) *HTMXRequest {
	tru := "true"
	isHTMXRequest := strings.TrimSpace(r.Header.Get(HXRequestHeader.String())) == tru
	if !isHTMXRequest {
		return nil
	}
	res := &HTMXRequest{
		Boosted:        strings.TrimSpace(r.Header.Get(BoostedRequest.String())) == tru,
		CurrentURL:     strings.TrimSpace(r.Header.Get(CurrentURLRequest.String())),
		HistoryRestore: strings.TrimSpace(r.Header.Get(HistoryRestoreRequest.String())) == tru,
		Prompt:         strings.TrimSpace(r.Header.Get(PromptRequest.String())),
		Target:         strings.TrimSpace(r.Header.Get(TargetRequest.String())),
		TriggerName:    strings.TrimSpace(r.Header.Get(TriggerNameRequest.String())),
		Trigger:        strings.TrimSpace(r.Header.Get(TriggerRequest.String())),
	}

	return res
}
