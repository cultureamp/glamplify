package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/stretchr/testify/assert"
)

func Test_New_TraceLogger(t *testing.T) {
	ctx := context.Background()

	log := newXrayLogger(ctx)
	assert.NotNil(t, log)
}

func Test_TraceLogger_RealWorld_Log(t *testing.T) {
	ctx := context.Background()

	log := newXrayLogger(ctx)
	log.Log(xraylog.LogLevelDebug, newPrintArgs("debug") )
	log.Log(xraylog.LogLevelInfo, newPrintArgs("info") )
	log.Log(xraylog.LogLevelWarn, newPrintArgs("warn") )
	log.Log(xraylog.LogLevelError, newPrintArgs("error") )
}


