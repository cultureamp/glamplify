package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	gcontext "github.com/cultureamp/glamplify/context"
	gerrors "github.com/go-errors/errors"
	"gotest.tools/assert"
)

var (
	ctx      context.Context
	rsFields gcontext.RequestScopedFields
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

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

func shutdown() {
	os.Unsetenv("PRODUCT")
	os.Unsetenv("APP")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_ACCOUNT_ID")
}

func Test_New(t *testing.T) {
	logger := New(rsFields)
	assert.Assert(t, logger != nil, logger)
}

func Test_NewWithContext(t *testing.T) {
	logger := NewFromCtx(ctx)
	assert.Assert(t, logger != nil, logger)

	rsFields, ok1 := gcontext.GetRequestScopedFields(ctx)

	assert.Assert(t, ok1, ok1)
	assert.Assert(t, rsFields.TraceID == "1-2-3", rsFields)
}

func Test_NewWithRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "*", nil)

	req1 := req.WithContext(ctx)
	logger := NewFromRequest(req1)
	assert.Assert(t, logger != nil, logger)

	rsFields, ok1 := gcontext.GetRequestScopedFields(req1.Context())

	assert.Assert(t, ok1, ok1)
	assert.Assert(t, rsFields.TraceID == "1-2-3", rsFields)
}

func Test_Log_IsEnabled(t *testing.T) {
	logger := NewFromCtx(ctx)
	assert.Assert(t, logger != nil, logger)

	assert.Assert(t, logger.IsEnabled(DebugSev), logger.IsEnabled(DebugSev))

	os.Setenv("LOG_LEVEL", InfoSev)
	defer os.Unsetenv("LOG_LEVEL")
	sevLevel := newSystemLogLevel()

	level := sevLevel.stringToLevel(DebugSev)
	assert.Assert(t, !sevLevel.shouldLog(level), sevLevel.shouldLog(level))
	level = sevLevel.stringToLevel(InfoSev)
	assert.Assert(t, sevLevel.shouldLog(level), sevLevel.shouldLog(level))
}

func Test_Log_Global_Scope(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Event( "detail_event").Fields(Fields{
		AppNameEnv: "app_name",
		AppFarmEnv: "app_farm",
	}).Debug("debug")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "detail_event")
	assertContainsString(t, msg, "severity", "DEBUG")
	assertContainsString(t, msg, "app", "app_name")
	assertContainsString(t, msg, "farm", "app_farm")
}

func Test_Log_Debug(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Debug( "detail_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "detail_event")
	assertContainsString(t, msg, "severity", "DEBUG")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_DebugWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "detail_event")
	assertContainsString(t, msg, "severity", "DEBUG")
	assertContainsString(t, msg, "string", "hello")
	assertContainsInt(t, msg, "int", 123)
	assertContainsFloat(t, msg, "float", 42.48)
	assertContainsString(t, msg, "string2", "hello world")
	assertContainsString(t, msg, "string3_space", "world")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
	assertScopeContainsSubDoc(t, msg, "properties")
}

func Test_Log_Info(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Info("info_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_InfoWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Info("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "string", "hello")
	assertContainsInt(t, msg, "int", 123)
	assertContainsFloat(t, msg, "float", 42.48)
	assertContainsString(t, msg, "string2", "hello world")
	assertContainsString(t, msg, "string3_space", "world")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
	assertScopeContainsSubDoc(t, msg, "properties")
}

func Test_Log_Warn(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Warn("warn_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "warn_event")
	assertContainsString(t, msg, "severity", "WARN")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_WarnWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Warn("warn_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "warn_event")
	assertContainsString(t, msg, "severity", "WARN")
	assertContainsString(t, msg, "string", "hello")
	assertContainsInt(t, msg, "int", 123)
	assertContainsFloat(t, msg, "float", 42.48)
	assertContainsString(t, msg, "string2", "hello world")
	assertContainsString(t, msg, "string3_space", "world")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
	assertScopeContainsSubDoc(t, msg, "properties")
}

func Test_Log_Error(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Error("error event", errors.New("something went wrong"))

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "error_event")
	assertContainsString(t, msg, "severity", "ERROR")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
	assertScopeContainsSubDoc(t, msg, "exception")
	assertContainsString(t, msg, "error", "something went wrong")
	//fmt.Println(msg)
}

func Test_Log_Error_StackTrace(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Error("error event", gerrors.New("with correct stack trace"))

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "error_event")
	assertContainsString(t, msg, "severity", "ERROR")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
	assertScopeContainsSubDoc(t, msg, "exception")
	assertContainsString(t, msg, "error", "with correct stack trace")
	//fmt.Println(msg)
}


func Test_Log_ErrorWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	logger.Error("error event", errors.New("something went wrong"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "error_event")
	assertContainsString(t, msg, "severity", "ERROR")
	assertContainsString(t, msg, "string", "hello")
	assertContainsInt(t, msg, "int", 123)
	assertContainsFloat(t, msg, "float", 42.48)
	assertContainsString(t, msg, "string2", "hello world")
	assertContainsString(t, msg, "string3_space", "world")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
	assertScopeContainsSubDoc(t, msg, "properties")
	assertScopeContainsSubDoc(t, msg, "exception")
	assertContainsString(t, msg, "error", "something went wrong")
}

func Test_Log_Fatal(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal_event")
			assertContainsString(t, msg, "severity", "FATAL")
			assertContainsString(t, msg, "trace_id", "1-2-3")
			assertContainsString(t, msg, "customer", "hooli")
			assertContainsString(t, msg, "user", "UserAggregateID-123")
			assertContainsString(t, msg, "product", "engagement")
			assertContainsString(t, msg, "app", "murmur")
			assertContainsString(t, msg, "app_version", "87.23.11")
			assertContainsString(t, msg, "aws_region", "us-west-02")
			assertContainsString(t, msg, "aws_account_id", "aws-account-123")
			assertScopeContainsSubDoc(t, msg, "exception")
			assertContainsString(t, msg, "error", "something fatal happened")
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
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal_event")
			assertContainsString(t, msg, "severity", "FATAL")
			assertContainsString(t, msg, "string", "hello")
			assertContainsInt(t, msg, "int", 123)
			assertContainsFloat(t, msg, "float", 42.48)
			assertContainsString(t, msg, "string2", "hello world")
			assertContainsString(t, msg, "string3_space", "world")
			assertContainsString(t, msg, "trace_id", "1-2-3")
			assertContainsString(t, msg, "customer", "hooli")
			assertContainsString(t, msg, "user", "UserAggregateID-123")
			assertContainsString(t, msg, "product", "engagement")
			assertContainsString(t, msg, "app", "murmur")
			assertContainsString(t, msg, "app_version", "87.23.11")
			assertContainsString(t, msg, "aws_region", "us-west-02")
			assertContainsString(t, msg, "aws_account_id", "aws-account-123")
			assertScopeContainsSubDoc(t, msg, "properties")
			assertScopeContainsSubDoc(t, msg, "exception")
			assertContainsString(t, msg, "error", "something fatal happened")

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

func Test_Log_Namespace(t *testing.T) {

	t1 := time.Now()
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	time.Sleep(123 * time.Millisecond)
	t2 := time.Now()
	d := t2.Sub(t1)

	logger.Error("error event", errors.New("something went wrong"), Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
		"reports_shared": Fields{
			"report":   "report1",
			"user":     "userid",
			TimeTaken: fmt.Sprintf("P%gS", d.Seconds()),
			TimeTakenMS: d.Milliseconds(),
		},
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "report", "report1")
	assertContainsString(t, msg, "user", "userid")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "hooli")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")

	assertScopeContainsSubDoc(t, msg, "reports_shared")
	assertScopeContainsSubDoc(t, msg, "properties")
}


func TestScope(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer, Fields{
		"requestID": 123,
	})

	logger.Debug("detail_event")
	msg := memBuffer.String()
	assertScopeContainsString(t, msg, "event", "detail_event")
	assertScopeContainsInt(t, msg, "request_id", 123)

	memBuffer.Reset()
	logger.Info("info_event")
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "info_event")
	assertScopeContainsInt(t, msg, "request_id", 123)

	memBuffer.Reset()
	logger.Warn("warn_event")
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "warn_event")
	assertScopeContainsInt(t, msg, "request_id", 123)

	memBuffer.Reset()
	logger.Error("error_event", errors.New("something went wrong"))
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "error_event")
	assertScopeContainsInt(t, msg, "request_id", 123)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal_event")
			assertContainsString(t, msg, "severity", "FATAL")
		}
	}()

	logger.Fatal("fatal_event", errors.New("something fatal happened")) // will call panic!
}

func TestScope_Overwrite(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer, Fields{
		"requestID": 123,
	})

	logger.Debug("detail_event", Fields {
		"requestID": 456,
	})
	msg := memBuffer.String()
	assertScopeContainsString(t, msg, "event", "detail_event")
	assertScopeContainsInt(t, msg, "request_id", 456)

	memBuffer.Reset()
	logger.Info("info_event", Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "info_event")
	assertScopeContainsInt(t, msg, "request_id", 456)

	memBuffer.Reset()
	logger.Warn("warn_event", Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "warn_event")
	assertScopeContainsInt(t, msg, "request_id", 456)

	memBuffer.Reset()
	logger.Error("error_event", errors.New("error"), Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "error_event")
	assertScopeContainsInt(t, msg, "request_id", 456)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertScopeContainsString(t, msg, "event", "fatal_event")
			assertScopeContainsString(t, msg, "severity", "FATAL")
			assertScopeContainsInt(t, msg, "request_id", 456)
		}
	}()

	// will call panic!
	logger.Fatal("fatal_event", errors.New("fatal"), Fields {
		"request_id": 456,
	})
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
	assert.Assert(t, logger != nil)

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

func Test_RealWorld_Durations(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	d := time.Millisecond * 456
	logger.Debug( "detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	}.Merge(NewDurationFields(d)))

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "detail_event")
	assertContainsString(t, msg, "time_taken", "P0.456S")
	assertContainsInt(t, msg, "time_taken_ms", 456)
}

func BenchmarkLogging(b *testing.B) {
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = ioutil.Discard
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

func assertContainsString(t *testing.T, log string, key string, val string) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertContainsInt(t *testing.T, log string, key string, val int) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":%v", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertContainsFloat(t *testing.T, log string, key string, val float32) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":%v", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertContainsSubDoc(t *testing.T, log string, key string, val string) {
	find := fmt.Sprintf("\"%s\":{\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)

}

func assertScopeContainsString(t *testing.T, log string, key string, val string) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertScopeContainsInt(t *testing.T, log string, key string, val int) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":%v", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertScopeContainsSubDoc(t *testing.T, log string, key string) {
	find := fmt.Sprintf("\"%s\":{", key)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}