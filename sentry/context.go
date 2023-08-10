package sentry

import (
	"context"
	"net/http"

	"github.com/go-errors/errors"
)

type key int

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
const (
	sentryContextKey key = iota
)

// FromRequest retrieves the current Application associated with the request, error is set appropriately
func FromRequest(_ http.ResponseWriter, r *http.Request) (*Application, error) {
	ctx := r.Context()
	return FromContext(ctx)
}

// FromContext gets the current Application from the given context
func FromContext(ctx context.Context) (*Application, error) {
	notify, ok := ctx.Value(sentryContextKey).(*Application)
	if ok && notify != nil {
		return notify, nil
	}

	return nil, errors.New("no sentry application in context")
}
