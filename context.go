package htmxtools

import "context"

type htmxContextKey string

var requestContextKey = htmxContextKey("htmxRequest")

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
