package bugsnag

import (
	"context"
	"errors"
	"net/http"
)

type key int

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
const (
	notifyContextKey     key = iota
)

// NotifyFromRequest retrieves the current Application associated with the request, error is set appropriately
func NotifyFromRequest(w http.ResponseWriter, r *http.Request) (*Application, error) {
	ctx := r.Context()
	return NotifyFromContext(ctx)
}

// NotifyFromContext gets the current Application from the given context
func NotifyFromContext(ctx context.Context) (*Application, error) {

	notify, ok := ctx.Value(notifyContextKey).(*Application)
	if ok && notify != nil {
		return notify, nil
	}

	return nil, errors.New("no notifier in context")
}

