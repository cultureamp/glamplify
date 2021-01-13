package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"

	gerrors "github.com/go-errors/errors"
)

// ErrorResponse response as per jsonapi.org - https://jsonapi.org/examples/#error-objects
type ErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`
}

// ErrorDetail represents a specific error as jsonapi.org - https://jsonapi.org/examples/#error-objects
type ErrorDetail struct {
	Code string `json:"code"`
	Title string `json:"title"`
	Detail string `json:"detail"`
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(code string, err error) *ErrorResponse {

	errors := &ErrorResponse{}
	errors = errors.AppendError(code, err)
	return errors
}

func (e *ErrorResponse) AppendError(code string, err error) *ErrorResponse {
	details := ErrorDetail{
		Code: code,
		Title: err.Error(),
		Detail: e.getLocation(err),
	}

	e.Errors = append(e.Errors, details)
	return e
}

// ToJSON returns json for an ErrorResponse
func (e ErrorResponse) ToJSON() string {
	b, err := json.Marshal(e)
	if err != nil {
		// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error#:~:text=From%20the%20docs%3A,result%20in%20an%20infinite%20recursion.
		// should not happen with a valid ErrorResponse
		panic(err)
	}
	return string(b)
}

func (e ErrorResponse) getLocation(err error) string {
	// is this error a https://github.com/go-errors/errors
	var goerr *gerrors.Error
	if errors.As(err, &goerr) {
		return e.getLocationFromErrorStack(goerr)
	}

	// skip 2 frames
	return e.getLocationFromCurrentRuntimeStack(2)
}

func (e ErrorResponse) getLocationFromErrorStack(err *gerrors.Error) string {

	callers := err.Callers()
	frames := runtime.CallersFrames(callers)

	// These frames are from the error, so we just need to get the first frame from that stack
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d:%s", frame.File, frame.Line, frame.Function)
}

func (e ErrorResponse) getLocationFromCurrentRuntimeStack(skip int) string {
	// We don't have a stack in the error, so walk the current runtime stack up
	// to the first caller that isn't glamplify....

	pc, file, line, ok := runtime.Caller(skip)
	for ok && strings.Contains(file, "glamplify") {
		skip++
		pc, file, line, ok = runtime.Caller(skip)
	}
	if !ok {
		return "unknown:0:unknown"
	}

	fn := runtime.FuncForPC(pc)
	methodName := fn.Name()
	return fmt.Sprintf("%s:%d:%s", file, line, methodName)
}