package sentry

import (
	"context"
	"errors"
	"net/http"
)

type key int

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
const (
	sentryContextKey key = iota
)

// SentryFromRequest retrieves the current Application associated with the request, error is set appropriately
func SentryFromRequest(w http.ResponseWriter, r *http.Request) (*Application, error) {
	ctx := r.Context()
	return SentryFromContext(ctx)
}

// SentryFromContext gets the current Application from the given context
func SentryFromContext(ctx context.Context) (*Application, error) {

	notify, ok := ctx.Value(sentryContextKey).(*Application)
	if ok && notify != nil {
		return notify, nil
	}

	return nil, errors.New("no sentry application in context")
}

