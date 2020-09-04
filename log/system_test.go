package log

import (
	"errors"
	"fmt"
	"testing"
	"time"

	gerrors "github.com/go-errors/errors"
	perrors "github.com/pkg/errors"
	"gotest.tools/assert"
)

func Test_HostName(t *testing.T) {

	df := newSystemValues()
	host := df.hostName()

	assert.Assert(t, host != "", host)
	assert.Assert(t, host != "<unknown>", host)
}

func Test_Default(t *testing.T) {
	df := newSystemValues()

	fields := df.getSystemValues(rsFields, "event_name", DebugSev)

	_, ok := fields[Time]
	assert.Assert(t, ok, "missing 'time' in default fields")
	_, ok = fields[Event]
	assert.Assert(t, ok, "missing 'event' in default fields")
	_, ok = fields[Resource]
	assert.Assert(t, ok, "missing 'resource' in default fields")
	_, ok = fields[Os]
	assert.Assert(t, ok, "missing 'os' in default fields")
	_, ok = fields[Severity]
	assert.Assert(t, ok, "missing 'severity' in default fields")

	_, ok = fields[TraceID]
	assert.Assert(t, ok, "missing 'trace_id' in default fields")
	_, ok = fields[Customer]
	assert.Assert(t, ok, "missing 'customer' in default fields")
	_, ok = fields[User]
	assert.Assert(t, ok, "missing 'user' in default fields")

	_, ok = fields[Product]
	assert.Assert(t, ok, "missing 'product' in default fields")
	_, ok = fields[App]
	assert.Assert(t, ok, "missing 'app' in default fields")
	_, ok = fields[AppVer]
	assert.Assert(t, ok, "missing 'app_ver' in default fields")
	_, ok = fields[AwsRegion]
	assert.Assert(t, ok, "missing 'region' in default fields")
}

func Test_ErrorDefault(t *testing.T) {
	df := newSystemValues()

	fields := df.getSystemValues(rsFields, "event_name", DebugSev)
	fields = df.getErrorValues(errors.New("test err"), fields)

	_, ok := fields[Exception]
	assert.Assert(t, ok, "missing 'exception' in default fields")
}

func Test_DurationAsIso8601(t *testing.T) {

	d := time.Millisecond * 456
	s := DurationAsISO8601(d)
	assert.Assert(t, s == "P0.456S", "was: %s", s)

	d = time.Millisecond * 1456
	s = DurationAsISO8601(d)
	assert.Assert(t, s == "P1.456S", "was: %s", s)
}

func Test_StackTrace(t *testing.T) {
	df := newSystemValues()

	stdStackFrame := df.getErrorStackTrace(errors.New("system error"))
	gStackFrame := df.getErrorStackTrace(gerrors.New("g error"))
	pStackFrame := df.getErrorStackTrace(perrors.New("p error"))

	assert.Assert(t, stdStackFrame != gStackFrame)
	assert.Assert(t, stdStackFrame != pStackFrame)
	fmt.Println("------ Standard Error Stack ------")
	fmt.Println(stdStackFrame)
	fmt.Println("------ go-errors ------")
	fmt.Println(gStackFrame)
	fmt.Println("------ pkg-errors ------")
	fmt.Println(pStackFrame)
}

func Test_CurrentStack(t *testing.T) {
	df := newSystemValues()

	stack0 := df.getCurrentStack(0)
	assert.Assert(t, stack0 != "")
	fmt.Println(stack0)

	stack1 := df.getCurrentStack(1)
	assert.Assert(t, stack1 != "")
	fmt.Println("------")
	fmt.Println(stack1)
}

