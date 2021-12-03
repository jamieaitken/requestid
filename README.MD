# RequestID 🔎
[![CircleCI](https://circleci.com/gh/jamieaitken/requestid/tree/main.svg?style=svg)](https://circleci.com/gh/jamieaitken/requestid/tree/main)
[![codecov](https://codecov.io/gh/jamieaitken/requestid/branch/main/graph/badge.svg?token=E7EV6JBPLW)](https://codecov.io/gh/jamieaitken/requestid)

This is a very simple piece of middleware for adding request/correlation IDs to the context of a http request.

By default, this module will add a key of `request-id` to context, however this can be overridden as shown below. 

## How to use

By default, you can just use the following:

```go
tracer := requestid.New()

router := new(mux.Router)
router.HandleFunc("/v1/get", tracer.Trace(HandlerFuncFoo))
router.HandleFunc("/v2/get", tracer.Trace(HandlerFuncBar))
```

If you wish to use another key, do the following

```go
const (
    CustomKey requestid.Key = "testKey"
)

tracer := requestid.New(requestid.WithTracerKey(CustomKey))

router := new(mux.Router)
router.HandleFunc("/v1/get", tracer.Trace(HandlerFuncFoo))
router.HandleFunc("/v2/get", tracer.Trace(HandlerFuncBar))
```

If you wish to use your own function, do the following

```go
tracer := requestid.New(requestid.WithTracerFunc(CustomFunc))

router := new(mux.Router)
router.HandleFunc("/v1/get", tracer.Trace(HandlerFuncFoo))
router.HandleFunc("/v2/get", tracer.Trace(HandlerFuncBar))

var (
    CustomKey requestid.Key = "testKey"
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
```

