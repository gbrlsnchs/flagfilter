package flagfilter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gbrlsnchs/flagfilter"
	"github.com/urfave/negroni"
)

func ExampleFilter() {
	r := httptest.NewRequest(http.MethodPost, "/example/test", nil)
	f := flagfilter.New(&flagfilter.Route{Path: `/example(/\w+)?`})
	allowed, err := f.Allowed(r)

	if err != nil {
		// handle error
	}

	fmt.Println(allowed)
	// Output: true
}

func ExampleFilter_context(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err, ok := r.Context().Value(flagfilter.CtxErrName).(error)

	if !ok {
		// handle type assertion error
	}

	if err != nil {
		// handle error
	}

	allowed, ok := r.Context().Value(flagfilter.CtxName).(bool)

	if !ok {
		// handle type assertion error
		return
	}

	if !allowed {
		// check auth and stuff
	}
}

func ExampleFilter_negroni(middlewareBefore, middlewareAfter negroni.Handler) {
	f := flagfilter.New(
		&flagfilter.Route{Path: "/example", Method: http.MethodGet},
		&flagfilter.Route{Path: "/auth", Method: http.MethodPost},
	)

	n := negroni.New(
		middlewareBefore,
		negroni.HandlerFunc(f.MiddlewareWithNext),
		middlewareAfter,
	)

	http.Handle("/api", n)
}
