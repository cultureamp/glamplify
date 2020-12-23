package log

import (
	"errors"
	"fmt"
	"testing"
	"time"

	gerrors "github.com/go-errors/errors"
	perrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_HostName(t *testing.T) {

	df := newSystemValues()
	host := df.hostName()

	assert.NotEmpty(t, host)
	assert.NotEqual(t, "<unknown>", host)
}

func Test_Default(t *testing.T) {
	df := newSystemValues()

	fields := df.getSystemValues(rsFields, nil, "event_name", DebugSev)

	_, ok := fields[Time]
	assert.True(t, ok)
	_, ok = fields[Event]
	assert.True(t, ok)
	_, ok = fields[Resource]
	assert.True(t, ok)
	_, ok = fields[Os]
	assert.True(t, ok)
	_, ok = fields[Severity]
	assert.True(t, ok)

	_, ok = fields[TraceID]
	assert.True(t, ok)
	_, ok = fields[Customer]
	assert.True(t, ok)
	_, ok = fields[User]
	assert.True(t, ok)

	_, ok = fields[Product]
	assert.True(t, ok)
	_, ok = fields[App]
	assert.True(t, ok)
	_, ok = fields[AppVer]
	assert.True(t, ok)
	_, ok = fields[Farm]
	assert.True(t, ok)
	_, ok = fields[AwsRegion]
	assert.True(t, ok)
}

func Test_ErrorDefault(t *testing.T) {
	df := newSystemValues()

	fields := df.getSystemValues(rsFields, nil, "event_name", DebugSev)
	fields = df.getErrorValues(errors.New("test err"), fields)

	_, ok := fields[Exception]
	assert.True(t, ok)
	_, ok = fields[Loc]
	assert.True(t, ok)
}

func Test_DurationAsIso8601(t *testing.T) {

	d := time.Millisecond * 456
	s := DurationAsISO8601(d)
	assert.Equal(t, "P0.456S", s)

	d = time.Millisecond * 1456
	s = DurationAsISO8601(d)
	assert.Equal(t, "P1.456S", s)
}

func Test_StackTrace(t *testing.T) {
	df := newSystemValues()

	stdStackFrame := df.getErrorStackTrace(errors.New("system error"))
	gStackFrame := df.getErrorStackTrace(gerrors.New("g error"))
	pStackFrame := df.getErrorStackTrace(perrors.New("p error"))

	fmt.Println("------ Standard Error Stack ------")
	fmt.Println(stdStackFrame)
	fmt.Println("------ go-errors ------")
	fmt.Println(gStackFrame)
	fmt.Println("------ pkg-errors ------")
	fmt.Println(pStackFrame)
	assert.NotEqual(t, stdStackFrame, gStackFrame)
	assert.NotEqual(t, stdStackFrame, pStackFrame)
}

func Test_CurrentStack(t *testing.T) {
	df := newSystemValues()

	stack0 := df.getCurrentStack(0)
	assert.NotEmpty(t, stack0)
	fmt.Println(stack0)

	stack1 := df.getCurrentStack(1)
	assert.NotEmpty(t, stack1)
	fmt.Println("------")
	fmt.Println(stack1)
}

