package flagfilter

import (
	"context"
	"net/http"
)

// Filter is responsible for flagging a route as allowed to continue its
// flow without authorization.
type Filter interface {
	// Allowed returns whether a route needs authorization.
	Allowed(*http.Request) (bool, error)
	// MiddlewareWithNext is a helper for the HTTP middleware pattern.
	MiddlewareWithNext(http.ResponseWriter, *http.Request, http.HandlerFunc)
	// Prefix sets a prefix for all routes that match the incoming request.
	Prefix(string) Filter
}

// filter holds the allowed routes and whether they were prefixed.
type filter struct {
	allowed  []*Route
	prefixed bool
}

// New creates a new Filter that receives
// a list of routes that don't need authorization.
func New(allowed ...*Route) Filter {
	return &filter{allowed: allowed}
}

// Allowed iterates over the filter's allowed routes and
// check for a match. If there's any, it returns early.
func (f *filter) Allowed(req *http.Request) (bool, error) {
	for _, r := range f.allowed {
		match, err := r.match(req)

		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	return false, nil
}

// MiddlewareWithNext will store the filter's flag inside the incoming request's context.
// If there's an error, it will store an error instead.
//
// Properties are named using CtxName and CtxErrName values.
func (f *filter) MiddlewareWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	allowed, err := f.Allowed(r)

	if err != nil {
		ctx := context.WithValue(r.Context(), CtxErrName, err)
		next(w, r.WithContext(ctx))

		return
	}

	ctx := context.WithValue(r.Context(), CtxName, allowed)

	next(w, r.WithContext(ctx))
}

// Prefix checks whether the routes were already prefixed,
// or if they're not an empty string. Then, it prefixes all
// routes according to the parameter passed.
func (f *filter) Prefix(prefix string) Filter {
	if f.prefixed || prefix == "" {
		return f
	}

	if len(f.allowed) > 0 {
		for _, r := range f.allowed {
			r.prefix(prefix)
		}

		f.prefixed = true
	}

	return f
}
