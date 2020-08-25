package newrelic

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/cultureamp/glamplify/log"
	"github.com/newrelic/go-agent/_integrations/nrlambda"
)

type lambdaHandler struct {
	impl lambda.Handler

	app Application

	functionName    string
	functionVersion string
	logGroupName    string
	logStreamName   string
	memoryLimitInMB int
}

// Invoke - Invoke API operation for AWS Lambda.
func (handler *lambdaHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	// payload => The JSON that you want to provide to your Lambda function as input.

	// Use handlertrace instead of wrapping the Invoke?
	// - https://godoc.org/github.com/aws/aws-lambda-go/lambda/handlertrace
	// - https://github.com/aws/aws-lambda-go/blob/master/lambda/handler_test.go

	var tin interface{}
	if handler.app.conf.Logging {
		json.Unmarshal(payload, &tin)

		handler.app.log("Begin Invoke", log.Fields{
			"function":        handler.functionName,
			"version":         handler.functionVersion,
			"logGroup":        handler.logGroupName,
			"logStream":       handler.logStreamName,
			"memoryLimitInMB": handler.memoryLimitInMB,
			"payload":         tin,
		})
	}

	// Add the CA handler to the ctx so that we can get it later inside "Invoke"
	ctx = handler.addToContext(ctx)
	// Add the CA application to the ctx so that we can get it later inside "Invoke"
	ctx = handler.app.addToContext(ctx)

	result, err := handler.impl.Invoke(ctx, payload)

	var tout interface{}
	if handler.app.conf.Logging {
		json.Unmarshal(result, &tout)
		handler.app.log("Invoke Ended", log.Fields{
			"function":        handler.functionName,
			"version":         handler.functionVersion,
			"logGroup":        handler.logGroupName,
			"logStream":       handler.logStreamName,
			"memoryLimitInMB": handler.memoryLimitInMB,
			"result":          tout,
			"error":           err,
		})
	}

	return result, err
}

// Start should be used in place of lambda.Start use app.Start(handler)
func (app Application) Start(handler interface{}) {
	// 1. First wrap the handler with NewRelic
	nr := nrlambda.Wrap(handler, app.impl)
	// 2. Then wrap that with CultureAmp
	ca := app.wrapLambda(nr)
	// 3. Start the handler
	lambda.StartHandler(ca)
}

// Start should be used in place of lambda.Start use Start(handler, app)
func Start(handler interface{}, app *Application) {
	app.Start(handler)
}

// StartHandler should be used in place of lambda.StartHandler use app.StartHandler(handler)
func (app Application) StartHandler(handler lambda.Handler) {
	// 1. First wrap the handler with NewRelic
	nr := nrlambda.WrapHandler(handler, app.impl)
	// 2. Then wrap that with CultureAmp
	ca := app.wrapLambdaHandler(nr)
	// 3. Start the handler
	lambda.StartHandler(ca)
}

// StartHandler should be used in place of lambda.StartHandler use StartHandler(handler, app)
func StartHandler(handler lambda.Handler, app *Application) {
	app.StartHandler(handler)
}

func (app Application) wrapLambda(handler lambda.Handler) lambda.Handler {

	return &lambdaHandler{
		impl:            handler,
		app:             app,
		functionName:    lambdacontext.FunctionName,
		functionVersion: lambdacontext.FunctionVersion,
		logGroupName:    lambdacontext.LogGroupName,
		logStreamName:   lambdacontext.LogStreamName,
		memoryLimitInMB: lambdacontext.MemoryLimitInMB,
	}
}

func (app Application) wrapLambdaHandler(handler interface{}) lambda.Handler {
	return app.wrapLambda(lambda.NewHandler(handler))
}

func (handler *lambdaHandler) addToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, handlerContextKey, handler)
}
