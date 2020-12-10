package log

import (
	"context"
	"fmt"
	gcontext "github.com/cultureamp/glamplify/context"
	"os"
	"testing"
)

func setup() {
	ctx = context.Background()
	ctx = gcontext.AddRequestFields(ctx, gcontext.RequestScopedFields{
		TraceID:             "1-2-3",
		RequestID:           "7-8-9",
		CorrelationID:       "1-5-9",
		CustomerAggregateID: "hooli",
		UserAggregateID:     "UserAggregateID-123",
	})

	rsFields, _ = gcontext.GetRequestScopedFields(ctx)

	os.Setenv("PRODUCT", "engagement")
	os.Setenv("APP", "murmur")
	os.Setenv("APP_ENV", "dev")
	os.Setenv("APP_VERSION", "87.23.11")
	os.Setenv("AWS_REGION", "us-west-02")
	os.Setenv("AWS_ACCOUNT_ID", "aws-account-123")
}

func teardown() {
	os.Unsetenv("PRODUCT")
	os.Unsetenv("APP")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_ACCOUNT_ID")
}

func TestMain(m *testing.M) {
	setup()
	runExitCode := m.Run()
	teardown()

	// runExitCode 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if runExitCode == 0 && testing.CoverMode() != "" {

		coverageResult := testing.Coverage()

		// If we are less than 90% then fail the build
		if coverageResult < 0.85 {
			fmt.Printf("Tests passed but coverage failed: MUST BE >= 85%%, was %.2f\n", coverageResult*100)
			runExitCode = -1
		}
	}
	os.Exit(runExitCode)
}