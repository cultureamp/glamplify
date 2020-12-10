package aws

import (
	"context"
	"testing"

	aws "github.com/aws/aws-lambda-go/events"
	"gotest.tools/assert"
)

type parameters map[string]string

func Test_Invoke_Handler(t *testing.T) {
	request := &aws.ALBTargetGroupRequest{
		Path:                  "/",
		HTTPMethod:            "GET",
		Headers:               parameters{"Keep-Alive": "timeout=5, max=1000"},
		QueryStringParameters: nil,
		IsBase64Encoded:       false,
	}

	response, err := InvokeLambda(Input{
		Port:               98001,
		Payload:            request,
		RequestID:          "0",
		InvokedFunctionArn: "arn:aws:lambda:us-east-1:123497558138:function:golang-layer:alias",
	})

	assert.Assert(t, err != nil, err)
	assert.Assert(t, len(response) == 0, len(response))
}

func handler(ctx context.Context, request aws.ALBTargetGroupRequest) (aws.ALBTargetGroupResponse, error) {
	return aws.ALBTargetGroupResponse{StatusCode: 200}, nil
}