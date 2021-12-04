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
	fn              TracerFunc
	key             Key
	idGeneratorFunc IDGeneratorFunc
}

// TracerFunc is the function signature of a Custom function to be used.
type TracerFunc func(http.HandlerFunc) http.HandlerFunc

type IDGeneratorFunc func() string

// Option allows for users to define the behaviour of a Tracer.
type Option func(*Tracer)

// New instantiates a Tracer.
func New(opts ...Option) *Tracer {
	t := &Tracer{
		key:             DefaultTracingKey,
		idGeneratorFunc: generateID,
	}

	t.fn = t.addRequestID

	t.Apply(opts...)

	return t
}

// Apply allows for the Tracer to be altered after instantiation.
func (t *Tracer) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(t)
	}
}

// Trace calls the given function.
func (t *Tracer) Trace(next http.HandlerFunc) http.HandlerFunc {
	return t.fn(next)
}

func generateID() string {
	return uuid.New().String()
}

func (t *Tracer) addRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		_, ok := ctx.Value(t.key).(string)
		if !ok {
			ctx = context.WithValue(ctx, t.key, t.idGeneratorFunc())
		}

		req = req.WithContext(ctx)

		next(w, req)
	}
}
