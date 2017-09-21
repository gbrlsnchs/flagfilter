package flagfilter

import (
	"net/http"
	"regexp"
	"strings"
)

// Route represents a route from an incoming request.
type Route struct {
	// Path is a regexp that matches incoming request's URL,
	// i.e., "/test" or "/test(ing)?".
	Path string
	// Method is a HTTP request method.
	Method string
	// Queries is an array of query string parameters.
	Queries []string
}

// match checks whether a Route's path, method and query string parameters
// matches the incoming request's path, method and query string parameters.
func (r *Route) match(req *http.Request) (bool, error) {
	var err error
	pathMatch := r.Path == ""

	if !pathMatch {
		pathMatch, err = regexp.MatchString(r.Path, req.URL.Path)

		if err != nil {
			return false, err
		}
	}

	methodMatch := r.Method == "" || r.Method == req.Method
	queryMatch := len(r.Queries) == 0
	q := req.URL.Query()

	for _, query := range r.Queries {
		_, queryMatch = q[query]

		if queryMatch {
			break
		}
	}

	return pathMatch && methodMatch && queryMatch, nil
}

// prefix prefixes a Route's path.
func (r *Route) prefix(prefix string) {
	if strings.HasPrefix(r.Path, "^") {
		r.Path = r.Path[1:]
		prefix = "^" + prefix
	}

	r.Path = prefix + r.Path
}
