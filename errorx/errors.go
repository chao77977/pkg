package errorx

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
)

type XError struct {
	code uint64
	msg  string
}

func New(code uint64, msg string) XError {
	return XError{code, msg}
}

func (x XError) Code() uint64 {
	return x.code
}

func (x XError) Error() string {
	return fmt.Sprintf("<errcode: %d> %s", x.code, x.msg)
}

func (x XError) Is(err error) bool {
	causeErr := errors.Cause(err)
	if causeErr, ok := causeErr.(XError); ok {
		return x.code == causeErr.code
	}

	return false
}

// unrecoverableError represents an unrecoverable error
type unrecoverableError struct {
	error
}

func (e unrecoverableError) Unwrap() error {
	return e.error
}

func (unrecoverableError) Is(err error) bool {
	_, isUnrecoverable := err.(unrecoverableError)
	return isUnrecoverable
}

// Unrecoverable returns an unrecoverable error
func Unrecoverable(err error) error {
	return unrecoverableError{err}
}

// IsRecoverable checks if error is unrecoverable error
func IsRecoverable(err error) bool {
	_, ok := err.(unrecoverableError)
	return !ok
}

// Errors represents list of errors
type Errors struct {
	offset int
	errs   []error
}

func NewErrors(offset int) *Errors {
	return &Errors{offset: offset}
}

func (e *Errors) Add(err error) {
	e.errs = append(e.errs, err)
}

func (e *Errors) Combine(errs ...error) {
	for _, err := range errs {
		if err != nil {
			e.Add(err)
		}
	}
}

func (e Errors) Is(err error) bool {
	for _, x := range e.errs {
		if errors.Is(x, err) {
			return true
		}
	}

	return false
}

// Unwrap pops up the last error for compatibility with errors.Unwrap()
func (e Errors) Unwrap() error {
	l := len(e.errs)
	if l < 1 {
		return nil
	}

	return e.errs[l-1]
}

// WrappedErrors returns the list of errors
func (e Errors) WrappedErrors() []error {
	return e.errs
}

func (e Errors) Error() string {
	l := len(e.errs)
	errs := make([]string, l)

	for i := 0; i < l; i++ {
		errs[i] = fmt.Sprintf("#%d: %s", i+e.offset, e.errs[i].Error())
	}

	return fmt.Sprintf("all failures with index:\n%s", strings.Join(errs, "\n"))
}
