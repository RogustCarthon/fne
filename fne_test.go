package fne_test

import (
	"errors"
	"testing"

	"github.com/RogustCarthon/fne"
	"github.com/stretchr/testify/assert"
)

func init() {
	fne.Separator(" > ")
}

type CustomErr struct {
	Message string
}

func (c CustomErr) Error() string {
	return c.Message
}

func TestError(t *testing.T) {
	errCause1 := fne.Rootf("something failed")
	errCause := fne.Errorf("failed to do that", errCause1)
	err := fne.Errorf("failed to do this", errCause)
	assert.Equal(t, `fne/fne_test.go::fne_test.TestError::26 failed to do this > fne/fne_test.go::fne_test.TestError::25 failed to do that > fne/fne_test.go::fne_test.TestError::24 something failed`, err.Error())
}

func TestNew(t *testing.T) {
	var err = errors.New("some err")
	fneErr := fne.New(err)
	assert.ErrorIs(t, fneErr, err)
}

func TestRootf(t *testing.T) {
	err := fne.Rootf("some err, data: %s", "some data")
	assert.Error(t, err)
}

func TestWrap(t *testing.T) {
	err1 := errors.New("err 1")
	err2 := CustomErr{"something went wrong"}
	err := fne.Wrap(err2, err1)
	assert.ErrorIs(t, err, err1)
	assert.ErrorIs(t, err, err2)
	assert.NotErrorIs(t, err2, err)
	assert.NotErrorIs(t, err1, err)
	var customErr = &CustomErr{}
	assert.ErrorAs(t, err, customErr)
}

func TestErrorf(t *testing.T) {
	err1 := errors.New("err 1")
	err2 := fne.Errorf("failed to do this", err1)
	err := fne.Errorf("failed to do that", err2)
	assert.ErrorIs(t, err, err1)
	assert.ErrorIs(t, err, err2)
	assert.ErrorIs(t, err2, err1)
	assert.NotErrorIs(t, err2, err)
	assert.NotErrorIs(t, err1, err)
}

func TestJoin(t *testing.T) {
	err1 := errors.New("err 1")
	err2 := errors.New("err 2")
	err3 := CustomErr{"something went wrong"}
	err := fne.Join(err1, err2, err3)

	assert.ErrorIs(t, err, err1)
	assert.ErrorIs(t, err, err2)
	assert.ErrorIs(t, err, err3)
	var customErr = &CustomErr{}
	assert.ErrorAs(t, err, customErr)

	err = fne.Join(nil)
	assert.Nil(t, err)
}
