package newrelic

import (
	"context"
	"net/http"

	"github.com/cultureamp/glamplify/log"
	newrelic "github.com/newrelic/go-agent"
)

// Transaction is a wrapper over the underlying implementation
type Transaction struct {
	impl    newrelic.Transaction
	app     *Application
	name    string
	logging bool
	logger  *monitorLogger
}

// GetApplication gets the Application from the current Transaction
func (txn Transaction) GetApplication() *Application {
	return txn.app
}

// AddAttributes adds customer data (key-value) to a current transaction (ie. http web request)
func (txn Transaction) AddAttributes(fields log.Fields) error {

	var err error
	for k, v := range fields {
		err = txn.impl.AddAttribute(k, v)
		if err != nil {
			txn.logError("AddAtributes", err)
			return err
		}
	}
	return nil
}

func (txn Transaction) ReportError(err error) error {
	return txn.impl.NoticeError(err)
}

func (txn Transaction) ReportErrorDetails(msg string, class string, fields log.Fields) error {
	return txn.impl.NoticeError(newrelic.Error{
		Message:    msg,
		Class:      class,
		Attributes: fields,
	})
}

// End closes the current transaction
func (txn Transaction) End() {
	txn.impl.End()

	txn.log("Transaction Ended", log.Fields{
		"txnName": txn.name,
	})
}

// Header delegates to the wrapped response
func (txn Transaction) Header() http.Header {
	return txn.impl.Header()
}

// Write delegates to the wrapped response
func (txn Transaction) Write(bytes []byte) (int, error) {
	count, err := txn.impl.Write(bytes)
	txn.logError("Write", err)
	return count, err
}

// WriteHeader delegates to the wrapped response
func (txn Transaction) WriteHeader(statusCode int) {
	txn.impl.WriteHeader(statusCode)
}

func (txn *Transaction) addToHTTPContext(req *http.Request) *http.Request {
	ctx := txn.addToContext(req.Context())
	return req.WithContext(ctx)
}

func (txn *Transaction) addToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, txnContextKey, txn)
}

func (txn Transaction) logError(msg string, err error) {
	if err != nil && txn.logging {
		txn.logger.Error(msg, map[string]interface{}{
			"error": err,
		})
	}
}

func (txn Transaction) log(msg string, fields log.Fields) {
	if txn.logging {
		txn.logger.Debug(msg, fields)
	}
}
