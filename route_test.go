package flagfilter

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteMatch(t *testing.T) {
	tests := []*struct {
		req      *http.Request
		r        *Route
		expected bool
	}{
		// #0
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			r:        &Route{Path: "/test"},
			expected: true,
		},
		// #1
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			r:        &Route{Path: "/testing"},
			expected: false,
		},
		// #2
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			r:        &Route{Path: "/test", Method: http.MethodGet},
			expected: true,
		},
		// #3
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			r:        &Route{Path: "/testing", Method: http.MethodGet},
			expected: false,
		},
		// #4
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			r:        &Route{Path: "/test", Method: http.MethodPost},
			expected: false,
		},
		// #5
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			r:        &Route{Path: "/testing", Method: http.MethodPost},
			expected: false,
		},
		// #6
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			r:        &Route{Path: "/test(ing)?"},
			expected: true,
		},
		// #7
		{
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			r:        &Route{Path: "/t*ing"},
			expected: false,
		},
		// #8
		{
			req:      httptest.NewRequest(http.MethodGet, "/test?key=value&value=key", nil),
			r:        &Route{Path: "/test", Method: http.MethodGet, Queries: []string{"key", "foo", "bar"}},
			expected: true,
		},
		// #9
		{
			req:      httptest.NewRequest(http.MethodGet, "/test?key=value&value=key", nil),
			r:        &Route{Path: "/test", Method: http.MethodGet, Queries: []string{"foo", "bar"}},
			expected: false,
		},
		// #10
		{
			req:      httptest.NewRequest(http.MethodPost, "/test/testing", nil),
			r:        &Route{Path: "/testing"},
			expected: true,
		},
		// #11
		{
			req:      httptest.NewRequest(http.MethodPost, "/test/testing", nil),
			r:        &Route{Path: "^/testing$"},
			expected: false,
		},
	}

	a := assert.New(t)
	for i, test := range tests {
		index := strconv.Itoa(i)
		match, err := test.r.match(test.req)

		a.Nil(err, index)
		a.Exactly(test.expected, match, index)
	}
}
