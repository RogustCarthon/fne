// package fne provides functionality to better chain errors with their causes.
package fne

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const (
	SeparatorNewLine = "\n"
	SeparatorArrow   = " > "
)

var separator = SeparatorNewLine

// Separator sets the separator string between cause and effected error
// messages. This should be called during initialization and is not concurrency
// safe.
func Separator(s string) {
	separator = s
}

type err struct {
	error
	cause      error
	function   string
	file       string
	lineNumber int
}

// Error method provides the string form of the error-cause tree.
func (e err) Error() string {
	var msg string
	if e.cause != nil {
		msg = fmt.Sprintf("%s%s%s", e.error.Error(), separator, e.cause.Error())
	} else {
		msg = e.error.Error()
	}
	return fmt.Sprintf(
		"%s::%s::%d %s",
		e.file, e.function, e.lineNumber, msg,
	)
}

// Is checks whether the target error matches any of the errors in the err-cause
// chain. It returns true if a match was found and false otherwise.
func (e err) Is(target error) bool {
	return errors.Is(e.error, target)
}

func (e err) Unwrap() error {
	return e.cause
}

// As checks whether the target matches any error in the err-cause chain and
// assigns the value of the first matched error to target. It returns true if a
// match was found and false otherwise.
func (e err) As(target any) bool {
	return errors.As(e.error, target)
}

// link should always be called by a caller which is 1 level above it.
func link(effect, cause error) error {
	pc, file, lineNumber, _ := runtime.Caller(2)
	fileParts := strings.Split(file, "/")
	file = strings.Join(fileParts[len(fileParts)-2:], "/")
	fnName := runtime.FuncForPC(pc).Name()
	fnName = fnName[strings.LastIndex(fnName, "/")+1:]
	return err{
		error:      effect,
		cause:      cause,
		function:   fnName,
		file:       file,
		lineNumber: lineNumber,
	}
}

func Rootf(message string, args ...any) error {
	return link(fmt.Errorf(message, args...), nil)
}

func Wrap(err, cause error) error {
	return link(err, cause)
}

func New(err error) error {
	return link(err, nil)
}

func Errorf(message string, err error, args ...any) error {
	return link(fmt.Errorf(message, args...), err)
}

func Join(errs ...error) error {
	err := errors.Join(errs...)
	if err != nil {
		return link(err, nil)
	}
	return nil
}
