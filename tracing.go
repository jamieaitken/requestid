package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Key is a type alias for the key to be given in context.
type Key string

const (
	DefaultTracingKey Key = "request-id"
)

// Tracer provides the ability to add a request id to context.
type Tracer struct {
	fn  TracerFunc
	key Key
}

// TracerFunc is the function signature of a Custom function to be used.
type TracerFunc func(http.HandlerFunc) http.HandlerFunc

// Option allows for users to define the behaviour of a Tracer.
type Option func(*Tracer)

// New instantiates a Tracer.
func New(opts ...Option) *Tracer {
	t := &Tracer{
		key: DefaultTracingKey,
	}

	t.fn = t.addRequestID

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// WithTracerFunc allows for using a custom function instead.
func WithTracerFunc(tracerFunc TracerFunc) Option {
	return func(tracer *Tracer) {
		tracer.fn = tracerFunc
	}
}

// WithTracerKey allows for a custom key to be used.
func WithTracerKey(tracerKey Key) Option {
	return func(tracer *Tracer) {
		tracer.key = tracerKey
	}
}

// Trace calls the given function.
func (t *Tracer) Trace(next http.HandlerFunc) http.HandlerFunc {
	return t.fn(next)
}

func (t *Tracer) addRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		_, ok := ctx.Value(t.key).(string)
		if !ok {
			ctx = context.WithValue(ctx, t.key, uuid.New().String())
		}

		req = req.WithContext(ctx)

		next(w, req)
	}
}
