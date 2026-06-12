package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type SomeCustomError struct {
	Message string
}

func (e *SomeCustomError) Error() string {
	return e.Message
}

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	builder := strings.Builder{}
	for _, err := range e.errors {
		builder.WriteString(err.Error() + "\t* ")
	}
	return strings.TrimSuffix(strconv.Itoa(len(e.errors))+" errors occured:\n\t* "+builder.String(), "\t* ") + "\n"
}

func (e *MultiError) Unwrap() []error {
	return e.errors
}

func Append(err error, errs ...error) *MultiError {
	var multiErr *MultiError
	if err != nil {
		if e, ok := err.(*MultiError); ok {
			multiErr = e
		} else {
			multiErr = &MultiError{errors: []error{err}}
		}
	} else {
		multiErr = &MultiError{}
	}

	multiErr.errors = append(multiErr.errors, errs...)
	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	err = Append(err, &SomeCustomError{Message: "custom error"})
	expectedMessage = "3 errors occured:\n\t* error 1\t* error 2\t* custom error\n"
	err = Append(err, errors.New("error 3"))
	expectedMessage = "4 errors occured:\n\t* error 1\t* error 2\t* custom error\t* error 3\n"
	assert.EqualError(t, err, expectedMessage)
	var target *SomeCustomError
	ok := errors.As(err, &target)
	assert.True(t, ok)
	ok = errors.Is(err, &SomeCustomError{})
	assert.False(t, ok)
	ok = errors.Is(err, target)
	assert.True(t, ok)
}
