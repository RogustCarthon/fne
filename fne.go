// package fne provides functionality to better chain errors with their causes.
package fne

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type err struct {
	error
	cause      error
	function   string
	file       string
	lineNumber int
}

func (e err) Error() string {
	var msg string
	if e.cause != nil {
		msg = fmt.Sprintf("%s > [ %s ]", e.error.Error(), e.cause.Error())
	} else {
		msg = e.error.Error()
	}
	return fmt.Sprintf(
		"%s::%s::%d %s",
		e.file, e.function, e.lineNumber, msg,
	)
}

func (e err) Is(err error) bool {
	return errors.Is(e.error, err)
}

func (e err) Unwrap() error {
	return e.cause
}

func (e err) As(target any) bool {
	return errors.As(e.error, target) || errors.As(e.cause, target)
}

func link(e, cause error) error {
	pc, file, lineNumber, _ := runtime.Caller(2)
	fileParts := strings.Split(file, "/")
	file = strings.Join(fileParts[len(fileParts)-2:], "/")
	fnName := runtime.FuncForPC(pc).Name()
	fnName = fnName[strings.LastIndex(fnName, "/")+1:]
	return err{
		error:      e,
		cause:      cause,
		function:   fnName,
		file:       file,
		lineNumber: lineNumber,
	}
}

func Rootf(message string, args ...any) error {
	return link(fmt.Errorf(message, args...), nil)
}

func Wrap(e, cause error) error {
	return link(e, cause)
}

func New(e error) error {
	return link(e, nil)
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
