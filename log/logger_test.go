package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/cultureamp/glamplify/env"

	gcontext "github.com/cultureamp/glamplify/context"
	gerrors "github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

var (
	ctx      context.Context
	rsFields gcontext.RequestScopedFields
)

func Test_New(t *testing.T) {
	logger := New(rsFields)
	assert.NotNil(t, logger)
}

func Test_NewWithContext(t *testing.T) {
	logger := NewFromCtx(ctx)
	assert.NotNil(t, logger)

	rsFields, ok := gcontext.GetRequestScopedFields(ctx)

	assert.True(t, ok)
	assert.Equal(t, "1-2-3", rsFields.TraceID)
}

func Test_NewWithRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "*", nil)

	req1 := req.WithContext(ctx)
	logger := NewFromRequest(req1)
	assert.NotNil(t, logger)

	rsFields, ok := gcontext.GetRequestScopedFields(req1.Context())

	assert.True(t, ok)
	assert.Equal(t, "1-2-3", rsFields.TraceID)
}

func Test_Log_IsEnabled(t *testing.T) {
	writer := NewWriter(func(config *WriterConfig) {
		config.Level = WarnSev
	})
	logger := NewFromCtxWithCustomerWriter(ctx, writer)
	assert.NotNil(t, logger)

	assert.False(t, logger.IsEnabled(DebugSev))
	assert.False(t, logger.IsEnabled(InfoSev))
	assert.True(t, logger.IsEnabled(WarnSev))
	assert.True(t, logger.IsEnabled(ErrorSev))
}

func Test_Log_Global_Scope(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Event("detail_event").Fields(Fields{
		env.AppNameEnv: "app_name",
		env.AppFarmEnv: "app_farm",
	}).Debug("debug")

	json := memBuffer.String()
	assert.Contains(t, json, "\"event\":\"detail_event\"")
	assert.Contains(t, json, "\"severity\":\"DEBUG\"")
	assert.Contains(t, json, "\"app\":\"app_name\"")
	assert.Contains(t, json, "\"farm\":\"app_farm\"")
}

func Test_Log_Debug(t *testing.T) {

	logger := New(rsFields)
	json := logger.Debug("detail_event")

	assert.Contains(t, json, "\"event\":\"detail_event\"")
	assert.Contains(t, json, "\"severity\":\"DEBUG\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")
}

func Test_Log_DebugWithFields(t *testing.T) {

	logger := New(rsFields)
	json := logger.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	assert.Contains(t, json, "\"event\":\"detail_event\"")
	assert.Contains(t, json, "\"severity\":\"DEBUG\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

	assert.Contains(t, json, "\"string\":\"hello\"")
	assert.Contains(t, json, "\"int\":123")
	assert.Contains(t, json, "\"float\":42.48")
	assert.Contains(t, json, "\"string2\":\"hello world\"")
	assert.Contains(t, json, "\"string3_space\":\"world\"")
}

func Test_Log_Info(t *testing.T) {

	logger := New(rsFields)
	json := logger.Info("info_event")

	assert.Contains(t, json, "\"event\":\"info_event\"")
	assert.Contains(t, json, "\"severity\":\"INFO\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")
}

func Test_Log_InfoWithFields(t *testing.T) {

	logger := New(rsFields)
	json := logger.Info("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	assert.Contains(t, json, "\"event\":\"info_event\"")
	assert.Contains(t, json, "\"severity\":\"INFO\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

	assert.Contains(t, json, "\"string\":\"hello\"")
	assert.Contains(t, json, "\"int\":123")
	assert.Contains(t, json, "\"float\":42.48")
	assert.Contains(t, json, "\"string2\":\"hello world\"")
	assert.Contains(t, json, "\"string3_space\":\"world\"")
}

func Test_Log_Warn(t *testing.T) {

	logger := New(rsFields)
	json := logger.Warn("warn_event")

	assert.Contains(t, json, "\"event\":\"warn_event\"")
	assert.Contains(t, json, "\"severity\":\"WARN\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")
}

func Test_Log_WarnWithFields(t *testing.T) {

	logger := New(rsFields)
	json := logger.Warn("warn_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	assert.Contains(t, json, "\"event\":\"warn_event\"")
	assert.Contains(t, json, "\"severity\":\"WARN\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

	assert.Contains(t, json, "\"string\":\"hello\"")
	assert.Contains(t, json, "\"int\":123")
	assert.Contains(t, json, "\"float\":42.48")
	assert.Contains(t, json, "\"string2\":\"hello world\"")
	assert.Contains(t, json, "\"string3_space\":\"world\"")
}

func Test_Log_Error(t *testing.T) {

	logger := New(rsFields)
	json := logger.Error("error event", errors.New("something went wrong"))

	assert.Contains(t, json, "\"event\":\"error_event\"")
	assert.Contains(t, json, "\"severity\":\"ERROR\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

	assert.Contains(t, json, "\"exception\"")
	assert.Contains(t, json, "\"error\":\"something went wrong\"")
}

func Test_Log_Error_StackTrace(t *testing.T) {

	logger := New(rsFields)
	json := logger.Error("error event", gerrors.New("with correct stack trace"))

	assert.Contains(t, json, "\"event\":\"error_event\"")
	assert.Contains(t, json, "\"severity\":\"ERROR\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

	assert.Contains(t, json, "\"exception\"")
	assert.Contains(t, json, "\"error\":\"with correct stack trace\"")
}

func Test_Log_ErrorWithFields(t *testing.T) {

	logger := New(rsFields)
	json := logger.Error("error event", errors.New("something went wrong"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	assert.Contains(t, json, "\"event\":\"error_event\"")
	assert.Contains(t, json, "\"severity\":\"ERROR\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

	assert.Contains(t, json, "\"string\":\"hello\"")
	assert.Contains(t, json, "\"int\":123")
	assert.Contains(t, json, "\"float\":42.48")
	assert.Contains(t, json, "\"string2\":\"hello world\"")
	assert.Contains(t, json, "\"string3_space\":\"world\"")

	assert.Contains(t, json, "\"exception\"")
	assert.Contains(t, json, "\"error\":\"something went wrong\"")
}

func Test_Log_Fatal(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	defer func() {
		if r := recover(); r != nil {
			json := memBuffer.String()
			assert.Contains(t, json, "\"event\":\"fatal_event\"")
			assert.Contains(t, json, "\"severity\":\"FATAL\"")
			assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
			assert.Contains(t, json, "\"customer\":\"hooli\"")
			assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
			assert.Contains(t, json, "\"product\":\"engagement\"")
			assert.Contains(t, json, "\"app\":\"murmur\"")
			assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
			assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
			assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

			assert.Contains(t, json, "\"exception\"")
			assert.Contains(t, json, "\"error\":\"something fatal happened\"")
		}
	}()

	logger.Fatal("fatal event", errors.New("something fatal happened")) // will call panic!
}

func Test_Log_FatalWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	defer func() {
		if r := recover(); r != nil {
			json := memBuffer.String()
			assert.Contains(t, json, "\"event\":\"fatal_event\"")
			assert.Contains(t, json, "\"severity\":\"FATAL\"")
			assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
			assert.Contains(t, json, "\"customer\":\"hooli\"")
			assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
			assert.Contains(t, json, "\"product\":\"engagement\"")
			assert.Contains(t, json, "\"app\":\"murmur\"")
			assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
			assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
			assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

			assert.Contains(t, json, "\"string\":\"hello\"")
			assert.Contains(t, json, "\"int\":123")
			assert.Contains(t, json, "\"float\":42.48")
			assert.Contains(t, json, "\"string2\":\"hello world\"")
			assert.Contains(t, json, "\"string3_space\":\"world\"")

			assert.Contains(t, json, "\"exception\"")
			assert.Contains(t, json, "\"error\":\"something fatal happened\"")
		}
	}()

	logger.Fatal("fatal event", errors.New("something fatal happened"), Fields{ // this will call panic!
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}

func Test_Log_Audit(t *testing.T) {

	logger := New(rsFields)
	json := logger.Audit("audit_event")

	assert.Contains(t, json, "\"event\":\"audit_event\"")
	assert.Contains(t, json, "\"severity\":\"AUDIT\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")
}

func Test_Log_AuditWithFields(t *testing.T) {

	logger := New(rsFields)
	json := logger.Audit("audit_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	assert.Contains(t, json, "\"event\":\"audit_event\"")
	assert.Contains(t, json, "\"severity\":\"AUDIT\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

	assert.Contains(t, json, "\"string\":\"hello\"")
	assert.Contains(t, json, "\"int\":123")
	assert.Contains(t, json, "\"float\":42.48")
	assert.Contains(t, json, "\"string2\":\"hello world\"")
	assert.Contains(t, json, "\"string3_space\":\"world\"")
}

func Test_Log_Error_WithSubDocument(t *testing.T) {

	t1 := time.Now()
	logger := New(rsFields)

	time.Sleep(123 * time.Millisecond)
	t2 := time.Now()
	d := t2.Sub(t1)
	timeTaken := fmt.Sprintf("P%gS", d.Seconds())

	json := logger.Error("error event", errors.New("something went wrong"), Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
		"reports_shared": Fields{
			"report":    "report1",
			"user":      "userid",
			TimeTaken:   timeTaken,
			TimeTakenMS: d.Milliseconds(),
		},
	})

	assert.Contains(t, json, "\"event\":\"error_event\"")
	assert.Contains(t, json, "\"severity\":\"ERROR\"")
	assert.Contains(t, json, "\"trace_id\":\"1-2-3\"")
	assert.Contains(t, json, "\"customer\":\"hooli\"")
	assert.Contains(t, json, "\"user\":\"UserAggregateID-123\"")
	assert.Contains(t, json, "\"product\":\"engagement\"")
	assert.Contains(t, json, "\"app\":\"murmur\"")
	assert.Contains(t, json, "\"app_version\":\"87.23.11\"")
	assert.Contains(t, json, "\"aws_region\":\"us-west-02\"")
	assert.Contains(t, json, "\"aws_account_id\":\"aws-account-123\"")

	assert.Contains(t, json, "\"string\":\"hello\"")
	assert.Contains(t, json, "\"int\":123")
	assert.Contains(t, json, "\"float\":42.48")
	assert.Contains(t, json, "\"reports_shared\"")
	assert.Contains(t, json, "\"report\":\"report1\"")
	assert.Contains(t, json, "\"user\":\"userid\"")
	assert.Contains(t, json, fmt.Sprintf("\"time_taken\":\"%s\"", timeTaken))
	assert.Contains(t, json, fmt.Sprintf("\"time_taken_ms\":%d", d.Milliseconds()))
}

func TestScope(t *testing.T) {
	logger := New(rsFields, Fields{
		"requestID": 123,
	})

	json := logger.Debug("detail_event")
	assert.Contains(t, json, "\"event\":\"detail_event\"")
	assert.Contains(t, json, "\"request_id\":123")

	json = logger.Info("info_event")
	assert.Contains(t, json, "\"event\":\"info_event\"")
	assert.Contains(t, json, "\"request_id\":123")

	json = logger.Warn("warn_event")
	assert.Contains(t, json, "\"event\":\"warn_event\"")
	assert.Contains(t, json, "\"request_id\":123")

	json = logger.Error("error_event", errors.New("something went wrong"))
	assert.Contains(t, json, "\"event\":\"error_event\"")
	assert.Contains(t, json, "\"request_id\":123")

	json = logger.Audit("audit_event")
	assert.Contains(t, json, "\"event\":\"audit_event\"")
	assert.Contains(t, json, "\"request_id\":123")

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger = NewWitCustomWriter(rsFields, writer, Fields{
		"requestID": 123,
	})

	defer func() {
		if r := recover(); r != nil {
			json := memBuffer.String()
			assert.Contains(t, json, "\"event\":\"fatal_event\"")
			assert.Contains(t, json, "\"request_id\":123")
		}
	}()

	logger.Fatal("fatal_event", errors.New("something fatal happened")) // will call panic!
}

func TestScope_Overwrite(t *testing.T) {

	logger := New(rsFields, Fields{
		"requestID": 123,
	})

	json := logger.Debug("detail_event", Fields{
		"requestID": 456,
	})
	assert.Contains(t, json, "\"event\":\"detail_event\"")
	assert.Contains(t, json, "\"request_id\":456")

	json = logger.Info("info_event", Fields{
		"requestID": 456,
	})
	assert.Contains(t, json, "\"event\":\"info_event\"")
	assert.Contains(t, json, "\"request_id\":456")

	json = logger.Warn("warn_event", Fields{
		"requestID": 456,
	})
	assert.Contains(t, json, "\"event\":\"warn_event\"")
	assert.Contains(t, json, "\"request_id\":456")

	json = logger.Error("error_event", errors.New("error"), Fields{
		"requestID": 456,
	})
	assert.Contains(t, json, "\"event\":\"error_event\"")
	assert.Contains(t, json, "\"request_id\":456")

	json = logger.Audit("audit_event", Fields{
		"requestID": 456,
	})
	assert.Contains(t, json, "\"event\":\"audit_event\"")
	assert.Contains(t, json, "\"request_id\":456")

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger = NewWitCustomWriter(rsFields, writer, Fields{
		"requestID": 123,
	})

	defer func() {
		if r := recover(); r != nil {
			json := memBuffer.String()
			assert.Contains(t, json, "\"event\":\"fatal_event\"")
			assert.Contains(t, json, "\"request_id\":456")
		}
	}()

	// will call panic!
	logger.Fatal("fatal_event", errors.New("fatal"), Fields{
		"request_id": 456,
	})
}

func Test_Durations(t *testing.T) {

	logger := New(rsFields)

	d := time.Millisecond * 456
	json := logger.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	}.Merge(NewDurationFields(d)))

	assert.Contains(t, json, "\"event\":\"detail_event\"")
	assert.Contains(t, json, "\"time_taken\":\"P0.456S\"")
	assert.Contains(t, json, "\"time_taken_ms\":456")
}

func Test_RealWorld(t *testing.T) {
	logger := New(rsFields)

	// You should see these printed out, all correctly formatted.
	logger.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	Debug(rsFields, "detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Info("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	Info(rsFields, "info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Warn("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	Warn(rsFields, "info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	logger.Audit("audit_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	Audit(rsFields, "audit_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Error("error_event", errors.New("error"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	Error(rsFields, "error_event", errors.New("error"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	defer func() {
		recover()
	}()

	// this will call panic!
	logger.Fatal("fatal_event", errors.New("fatal"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	defer func() {
		recover()
	}()

	// this will call panic!
	Fatal(rsFields, "fatal_event", errors.New("fatal"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}

func Test_RealWorld_Combined(t *testing.T) {
	logger := New(rsFields)

	// multiple fields collections
	logger.Debug("detail_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
	Debug(rsFields, "detail_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	logger.Info("info_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
	Info(rsFields, "info_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	logger.Warn("warn_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
	Warn(rsFields, "warn_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	logger.Audit("audit_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
	Audit(rsFields, "audit_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	logger.Error("error_event", errors.New("error"), Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
	Error(rsFields, "error_event", errors.New("error"), Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	defer func() {
		recover()
	}()

	// this will call panic!
	logger.Fatal("fatal_event", errors.New("fatal"), Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	defer func() {
		recover()
	}()

	// this will call panic!
	Fatal(rsFields, "fatal_event", errors.New("fatal"), Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
}

func Test_RealWorld_Scope(t *testing.T) {

	logger := New(rsFields, Fields{"scopeID": 123})
	assert.NotNil(t, logger)

	logger.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Info("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Warn("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Audit("audit_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Error("error_event", errors.New("error"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	defer func() {
		recover()
	}()

	// this will call panic!
	logger.Fatal("fatal_event", errors.New("fatal"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}

func BenchmarkLogging(b *testing.B) {
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = io.Discard
	})
	logger := newLogger(rsFields, writer)

	fields := Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	}

	for n := 0; n < b.N; n++ {
		logger.Info("test details", fields)
	}
}
