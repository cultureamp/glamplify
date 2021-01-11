package main

import (
	"errors"
	"testing"

	gerrors "github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func Test_ErrorResponse_New(t *testing.T) {

	err := gerrors.New("test")
	er := NewErrorResponse("500", err)
	assert.NotNil(t, er)
	assert.Len(t, er.Errors, 1)

	s := er.ToJSON()
	assert.NotEmpty(t, s)
	//fmt.Println(s)
	assert.JSONEq(t, "{\"errors\":[{\"code\":\"500\",\"title\":\"test\",\"detail\":\"C:/src/go/src/github.com/cultureamp/glamplify/errors/errors_test.go:13:github.com/cultureamp/glamplify/errors.Test_ErrorResponse_New\"}]}", s)
}

func Test_ErrorResponse_Append(t *testing.T) {
	err := gerrors.New("first error")
	er := NewErrorResponse("500", err)
	assert.NotNil(t, er)
	assert.Len(t, er.Errors, 1)

	er = er.AppendError("404", gerrors.New("second error"))
	assert.Len(t, er.Errors, 2)

	s := er.ToJSON()
	assert.NotEmpty(t, s)
	//fmt.Println(s)
	assert.JSONEq(t, "{\"errors\":[{\"code\":\"500\",\"title\":\"first error\",\"detail\":\"C:/src/go/src/github.com/cultureamp/glamplify/errors/errors_test.go:25:github.com/cultureamp/glamplify/errors.Test_ErrorResponse_Append\"},{\"code\":\"404\",\"title\":\"second error\",\"detail\":\"C:/src/go/src/github.com/cultureamp/glamplify/errors/errors_test.go:30:github.com/cultureamp/glamplify/errors.Test_ErrorResponse_Append\"}]}", s)
}

func Test_ErrorResponse_Standard_Error(t *testing.T) {
	err := errors.New("standard error")
	er := NewErrorResponse("500", err)
	assert.NotNil(t, er)
	assert.Len(t, er.Errors, 1)

	s := er.ToJSON()
	assert.NotEmpty(t, s)
	//fmt.Println(s)
	assert.JSONEq(t, "{\"errors\":[{\"code\":\"500\",\"title\":\"standard error\",\"detail\":\"C:/Program Files/go/src/testing/testing.go:1123:testing.tRunner\"}]}", s)
}