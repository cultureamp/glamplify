package httplambda

import (
	"context"
	"fmt"
	"glamplify/server/datadog"
	"net/http"

	"github.com/apex/gateway/v2"
	"github.com/aws/aws-lambda-go/events"
)

// The Apex Gateway library exposes an implementation of the lambda.Handler
// interface, which is great for calling lambda.StartHandler(). If you
// use lambda.Start() instead, reflection is used to unmarshal the event
// and supply it to the method. Since Datadog only supports the reflection-based method,
// we need to use the Apex primitives and implement the expected function shape.
//
// See https://pkg.go.dev/github.com/aws/aws-lambda-go@v1.22.0/lambda#Start

// NewAdaptor creates a function suitable for passing to lambda.InvokeWithContext
// that wraps a standard http request.
func NewAdaptor(handler http.Handler) interface{} {
	return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		datadog.RecordRequestStart(ctx, event.RequestContext.HTTP.Method, event.RequestContext.HTTP.Path)

		if _, ok := event.Headers["accept"]; ok {
			fmt.Println("FIXME: overriding accept header")
			event.Headers["accept"] = "application/json"
		}

		fmt.Printf("\n**************\nAPIGatewayV2HTTPRequest: %+v\n**************\n", event.Body)
		r, err := gateway.NewRequest(ctx, event)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{}, err
		}

		// The Apex library overwrites an incoming request ID with the API
		// gateway request ID. If a request ID was supplied in a header, make
		// sure it takes precedence.
		if requestId, ok := event.Headers["x-request-id"]; ok {
			r.Header.Set("x-request-id", requestId)
		}

		w := gateway.NewResponse()
		handler.ServeHTTP(w, r)

		resp := w.End()

		datadog.RecordRequestComplete(ctx, resp.StatusCode)

		return resp, nil
	}
}
