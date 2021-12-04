package requestid_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jamieaitken/requestid"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name             string
		givenOpts        []requestid.Option
		givenHandlerFunc http.HandlerFunc
		expectedStatus   int
	}{
		{
			name:             "given request, expect default requestid function to be used",
			givenHandlerFunc: Get,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "given request, expect default requestid function to be used with given key",
			givenOpts:        []requestid.Option{requestid.WithTracerKey(CustomKey)},
			givenHandlerFunc: GetWithGivenKey,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "given request and custom func, expect custom func to be used",
			givenOpts:        []requestid.Option{requestid.WithTracerFunc(CustomFunc)},
			givenHandlerFunc: GetCustom,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "given request and custom id generator, expect static-key in context",
			givenOpts:        []requestid.Option{requestid.WithIDGenerator(staticID)},
			givenHandlerFunc: GetWithCustomIDGenerator,
			expectedStatus:   http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tracer := requestid.New(test.givenOpts...)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/v1/get", nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/get", tracer.Trace(test.givenHandlerFunc))
			router.ServeHTTP(rr, req)

			if !cmp.Equal(rr.Code, test.expectedStatus) {
				t.Fatalf(cmp.Diff(rr.Code, test.expectedStatus))
			}
		})
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(requestid.DefaultTracingKey).(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetWithGivenKey(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(CustomKey).(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

const (
	CustomKey requestid.Key = "testKey"
	StaticKey string        = "static-key"
)

func CustomFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		_, ok := ctx.Value(CustomKey).(string)
		if !ok {
			ctx = context.WithValue(ctx, CustomKey, uuid.New().String())
		}

		req = req.WithContext(ctx)

		next(w, req)
	}
}

func GetCustom(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(CustomKey).(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetWithCustomIDGenerator(w http.ResponseWriter, r *http.Request) {
	val, ok := r.Context().Value(requestid.DefaultTracingKey).(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if val != StaticKey {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func staticID() string {
	return StaticKey
}
