package aws

import (
	"encoding/json"
	"fmt"
	"net/rpc"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
	lc "github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/go-errors/errors"
)

const (
	functionInvokeRPC = "Function.Invoke"
)

// Input represents an AWS RPC call
type Input = struct {
	Port                  int
	Payload               interface{}
	ClientContext         *lc.ClientContext
	Deadline              *messages.InvokeRequest_Timestamp
	RequestID             string
	XAmznTraceID          string
	InvokedFunctionArn    string
	CognitoIdentityID     string
	CognitoIdentityPoolID string
}

//InvokeLambda a Go based lambda, passing the configured payload
func InvokeLambda(input Input) ([]byte, error) {
	request, err := createInvokeRequest(input)

	if err != nil {
		return nil, err
	}

	// 2. Open a TCP connection to the lambda
	client, err := rpc.Dial("tcp", fmt.Sprintf(":%d", input.Port))
	if err != nil {
		return nil, err
	}

	// 3. Issue an RPC request for the Function.Invoke method
	var response messages.InvokeResponse

	if err = client.Call(functionInvokeRPC, request, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}

	return response.Payload, nil
}

func createInvokeRequest(input Input) (*messages.InvokeRequest, error) {
	payloadEncoded, err := json.Marshal(input.Payload)
	if err != nil {
		return nil, err
	}

	var clientContextEncoded []byte
	if input.ClientContext != nil {
		b, err := json.Marshal(input.ClientContext)

		if err != nil {
			return nil, err
		}

		clientContextEncoded = b
	}

	Deadline := input.Deadline

	if Deadline == nil {
		t := time.Now().Add(5 * time.Second)
		Deadline = &messages.InvokeRequest_Timestamp{
			Seconds: t.Unix(),
			Nanos:   int64(t.Nanosecond()),
		}
	}

	return &messages.InvokeRequest{
		Payload:               payloadEncoded,
		RequestId:             input.RequestID,
		XAmznTraceId:          input.XAmznTraceID,
		Deadline:              *Deadline,
		InvokedFunctionArn:    input.InvokedFunctionArn,
		CognitoIdentityId:     input.CognitoIdentityID,
		CognitoIdentityPoolId: input.CognitoIdentityPoolID,
		ClientContext:         clientContextEncoded,
	}, nil
}
