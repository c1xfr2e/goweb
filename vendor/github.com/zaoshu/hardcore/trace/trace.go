package trace

import (
	"net/http"

	"github.com/zaoshu/hardcore/uuid"
	"golang.org/x/net/context"
)

// RequestIDKey key of request id in context
const RequestIDKey = "__zaoshu/hardcore/trace/key/request__"

// NewRequestID new request id
func NewRequestID() string {
	return uuid.New()
}

// GetRequestID get request id from context
func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDKey).(string)
	return id
}

// NewContext new context with request id
func NewContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// SetHTTPHeaderFromContext set http header
func SetHTTPHeaderFromContext(ctx context.Context, header http.Header) {
	id := GetRequestID(ctx)
	if len(id) == 0 {
		id = NewRequestID()
	}
	header.Set("X-Request-Id", id)
}
