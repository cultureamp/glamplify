package bugsnag

import (
	"context"
	"errors"
	"net/http"
)

type key int

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
const (
	bugsnagContextKey key = iota
)

// FromRequest retrieves the current Application associated with the request, error is set appropriately
func FromRequest(w http.ResponseWriter, r *http.Request) (*Application, error) {
	ctx := r.Context()
	return FromContext(ctx)
}

// FromContext gets the current Application from the given context
func FromContext(ctx context.Context) (*Application, error) {

	notify, ok := ctx.Value(bugsnagContextKey).(*Application)
	if ok && notify != nil {
		return notify, nil
	}

	return nil, errors.New("no bugsnag application in context")
}

