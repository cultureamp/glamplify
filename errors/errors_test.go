package errors

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

	// build has different paths, so check it has these "snippets" which shouldn't change
	assert.Contains(t, s, "{\"errors\":[{\"code\":\"500\",", s)
	assert.Contains(t, s, "\"title\":\"test\"", s)
	assert.Contains(t, s, "errors.Test_ErrorResponse_New\"}]}", s)
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

	// build has different paths, so check it has these "snippets" which shouldn't change
	assert.Contains(t, s, "\"code\":\"500\",", s)
	assert.Contains(t, s, "\"code\":\"404\",", s)
	assert.Contains(t, s, "\"title\":\"first error\"", s)
	assert.Contains(t, s, "\"title\":\"second error\"", s)
	assert.Contains(t, s, "errors.Test_ErrorResponse_Append\"}]}", s)
}

func Test_ErrorResponse_Standard_Error(t *testing.T) {
	err := errors.New("standard error")
	er := NewErrorResponse("500", err)
	assert.NotNil(t, er)
	assert.Len(t, er.Errors, 1)

	s := er.ToJSON()
	assert.NotEmpty(t, s)
	//fmt.Println(s)

	// build has different paths, so check it has these "snippets" which shouldn't change
	assert.Contains(t, s, "{\"errors\":[{\"code\":\"500\",", s)
	assert.Contains(t, s, "\"title\":\"standard error\"", s)
}