package aws

import (
	"context"
	"testing"

	aws "github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
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

	assert.NotNil(t, err)
	assert.Len(t, response, 0)
}

func handler(ctx context.Context, request aws.ALBTargetGroupRequest) (aws.ALBTargetGroupResponse, error) {
	return aws.ALBTargetGroupResponse{StatusCode: 200}, nil
}