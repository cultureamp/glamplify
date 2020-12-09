package newrelic

import (
	"context"
	"errors"
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

type key int

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
const (
	txnContextKey     key = iota
	appContextKey     key = iota
	handlerContextKey key = iota
)

// TxnFromRequest retrieves the current Transaction associated with the request, error is set appropriately
func TxnFromRequest(w http.ResponseWriter, r *http.Request) (*Transaction, error) {
	ctx := r.Context()
	return TxnFromContext(ctx)
}

// TxnFromContext gets the current Transaction from the given context
func TxnFromContext(ctx context.Context) (*Transaction, error) {

	// 1. First try and get the CA txn from the context. It will be there for HTTP wrapped methods,
	// but not for serverless/lambda ones
	txn, err := txnFromContext(ctx)
	if err == nil && txn != nil {
		return txn, nil
	}

	// 2. So likely a serverless/lambda call, so try and get the CA lambdaHandler so we can get the txnName & app
	// It should be there as we added it before calling "Invoke".
	handler, err := handlerFromContext(ctx)
	if err != nil || handler == nil {
		return nil, err
	}

	// 3. So likely a serverless/lambda call, so get the NR txn from the ctx
	txnName := handler.functionName
	impl := newrelic.FromContext(ctx)
	if impl != nil {
		// A bit yuck - we need to create a CA txn here after the fact because NR created one invisibly to us...
		txn = &Transaction{
			impl:    impl,
			app:     &handler.app,
			name:    txnName,
			logging: handler.app.conf.Logging,
			logger:  handler.app.conf.logger,
		}

		return txn, nil
	}

	// No transaction!
	err = errors.New("no transaction found")
	handler.app.logError("Call app.StartTransaction() to create a new transaction.", err)
	return nil, err
}

// TxnFromContext retireves the current Transation from the given context, error is set appropriately
func txnFromContext(ctx context.Context) (*Transaction, error) {
	txn, ok := ctx.Value(txnContextKey).(*Transaction)
	if ok && txn != nil {
		return txn, nil
	}

	err := errors.New("no transaction in context")
	return nil, err
}

// AppFromRequest retrives the current Application associated with the request, error is set appropriately
func AppFromRequest(w http.ResponseWriter, r *http.Request) (*Application, error) {
	ctx := r.Context()
	return AppFromContext(ctx)
}

// AppFromContext retireves the current Application from the given context, error is set appropriately
func AppFromContext(ctx context.Context) (*Application, error) {
	app, ok := ctx.Value(appContextKey).(*Application)
	if ok && app != nil {
		return app, nil
	}

	err := errors.New("no application in context")
	return nil, err
}

func handlerFromContext(ctx context.Context) (*lambdaHandler, error) {
	handler, ok := ctx.Value(handlerContextKey).(*lambdaHandler)
	if ok && handler != nil {
		return handler, nil
	}

	err := errors.New("no handler in context")
	return nil, err
}
