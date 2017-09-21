package flagfilter

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFilter(t *testing.T) {
	a := assert.New(t)
	tests := []*struct {
		req      *http.Request
		f        Filter
		expected bool
	}{
		// #0
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			f:        New(&Route{Path: "/test"}),
			expected: true,
		},
		// #1
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			f:        New(),
			expected: false,
		},
		// #2
		{
			req:      httptest.NewRequest(http.MethodPost, "/test/testing", nil),
			f:        New(&Route{Path: "^/testing$"}).Prefix("/test"),
			expected: true,
		},
		// #3
		{
			req:      httptest.NewRequest(http.MethodPost, "/test/testing", nil),
			f:        New(&Route{Path: "^/testing$"}).Prefix("/tester"),
			expected: false,
		},
	}

	for i, test := range tests {
		allowed, err := test.f.Allowed(test.req)
		index := strconv.Itoa(i)

		a.Nil(err, index)
		a.Exactly(test.expected, allowed, index)

		w := httptest.NewRecorder()

		test.f.MiddlewareWithNext(w, test.req, func(w http.ResponseWriter, r *http.Request) {
			err := r.Context().Value(CtxErrName)

			a.Nil(err, index)

			allowed, ok := r.Context().Value(CtxName).(bool)

			a.True(ok, index)
			a.Exactly(test.expected, allowed, index)
		})
	}
}
