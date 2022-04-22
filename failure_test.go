package failure_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"

	"github.com/rsb/failure"
	"github.com/stretchr/testify/assert"
)

func TestBadRequest(t *testing.T) {
	msg := "some message"
	err := failure.BadRequest(msg)
	assert.Error(t, err, "failure.BadRequest is expected to return an error")
	assert.True(t, failure.IsBadRequest(err))
	assert.Contains(t, err.Error(), failure.BadRequestMsg)

	assert.False(t, failure.IsSystem(err))

}

func TestToBadRequest(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToBadRequest(e, "user messed up")
	assert.Error(t, err, "failure.ToBadRequest is expected to return an error")
	assert.True(t, failure.IsBadRequest(err))

	expected := "user messed up: api specific msg: bad request"
	assert.Equal(t, err.Error(), expected)
}

func TestServer(t *testing.T) {
	msg := "some message"
	err := failure.Server(msg)
	assert.Error(t, err, "failure.Server is expected to return an error")
	assert.True(t, failure.IsServer(err))
	assert.Contains(t, err.Error(), failure.ServerMsg)

	// They are conceptually the same but are not equal
	assert.False(t, failure.IsSystem(err))

}

func TestToServer(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToServer(e, "api x failed")
	assert.Error(t, err, "failure.ToServer is expected to return an error")
	assert.True(t, failure.IsServer(err))

	expected := "api x failed: api specific msg: server failure"
	assert.Equal(t, err.Error(), expected)
}

func TestSystem(t *testing.T) {
	msg := "some message"
	err := failure.System(msg)
	assert.Error(t, err, "failure.System is expected to return an error")
	assert.True(t, failure.IsSystem(err))

	assert.False(t, failure.IsServer(err))

	assert.Contains(t, err.Error(), failure.SystemMsg)
}

func TestToSystem(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToSystem(e, "api x failed")
	assert.Error(t, err, "failure.ToSystem is expected to return an error")
	assert.True(t, failure.IsSystem(err))

	expected := "api x failed: api specific msg: system failure"
	assert.Equal(t, err.Error(), expected)
}

func TestShutdown(t *testing.T) {
	msg := "some message"
	err := failure.Shutdown(msg)
	assert.Error(t, err, "failure.Shutdown is expected to return an error")
	assert.True(t, failure.IsShutdown(err))

	assert.False(t, failure.IsServer(err))

	assert.Contains(t, err.Error(), failure.ShutdownMsg)
}

func TestToShutdown(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToShutdown(e, "api x failed")
	assert.Error(t, err, "failure.ToShutdown is expected to return an error")
	assert.True(t, failure.IsShutdown(err))

	expected := "api x failed: api specific msg: system shutdown failure"
	assert.Equal(t, err.Error(), expected)
}

func TestConfig(t *testing.T) {
	msg := "some message"
	err := failure.Config(msg)
	assert.Error(t, err, "failure.Config is expected to return an error")
	assert.True(t, failure.IsConfig(err))

	assert.False(t, failure.IsServer(err))
	assert.False(t, failure.IsSystem(err))

	assert.Contains(t, err.Error(), failure.ConfigMsg)
}

func TestInvalidParam(t *testing.T) {
	msg := "some message"
	err := failure.InvalidParam(msg)
	assert.Error(t, err, "failure.InvalidParam is expected to return an error")
	assert.True(t, failure.IsInvalidParam(err))

	assert.False(t, failure.IsServer(err))
	assert.False(t, failure.IsSystem(err))

	assert.Contains(t, err.Error(), failure.InvalidParamMsg)
}

func TestToInvalidParam(t *testing.T) {
	msg := "env var error msg"
	e := errors.New(msg)

	err := failure.ToInvalidParam(e, "some param failed")
	assert.Error(t, err, "failure.ToConfig is expected to return an error")
	assert.True(t, failure.IsInvalidParam(err))

	expected := fmt.Sprintf("some param failed: env var error msg: %s", failure.InvalidParamMsg)
	assert.Equal(t, expected, err.Error())
}

func TestToConfig(t *testing.T) {
	msg := "env var error msg"
	e := errors.New(msg)

	err := failure.ToConfig(e, "some env failed")
	assert.Error(t, err, "failure.ToConfig is expected to return an error")
	assert.True(t, failure.IsConfig(err))

	expected := fmt.Sprintf("some env failed: env var error msg: %s", failure.ConfigMsg)
	assert.Equal(t, expected, err.Error())
}

func TestNotFound(t *testing.T) {
	msg := "some message"
	err := failure.NotFound(msg)
	assert.Error(t, err, "failure.NotFound is expected to return an error")
	assert.True(t, failure.IsNotFound(err))

	assert.False(t, failure.IsServer(err))
	assert.False(t, failure.IsSystem(err))
	assert.False(t, failure.IsConfig(err))

	assert.Contains(t, err.Error(), failure.NotFoundMsg)
}

func TestToNotFound(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToNotFound(e, "api resource does not exist")
	assert.Error(t, err, "failure.ToNotFound is expected to return an error")
	assert.True(t, failure.IsNotFound(err))

	expected := fmt.Sprintf("api resource does not exist: api specific msg: %s", failure.NotFoundMsg)
	assert.Equal(t, err.Error(), expected)
}

func TestIgnore(t *testing.T) {
	msg := "some message"
	err := failure.Ignore(msg)
	assert.Error(t, err, "failure.Ignore is expected to return an error")

	assert.Contains(t, err.Error(), failure.IgnoreMsg)
}

func TestIsIgnore(t *testing.T) {
	err := failure.Ignore("some message")
	assert.True(t, failure.IsIgnore(err))

	assert.False(t, failure.IsServer(err))
	assert.False(t, failure.IsSystem(err))
	assert.False(t, failure.IsConfig(err))
	assert.False(t, failure.IsNotFound(err))
	assert.False(t, failure.IsValidation(err))
}

func TestToIgnore(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToIgnore(e, "api error should be ignored")
	assert.Error(t, err, "failure.ToIgnore is expected to return an error")
	assert.True(t, failure.IsIgnore(err))

	expected := fmt.Sprintf("api error should be ignored: api specific msg: %s", failure.IgnoreMsg)
	assert.Equal(t, err.Error(), expected)
}

func TestValidation(t *testing.T) {
	msg := "some message"
	err := failure.Validation(msg)
	assert.Error(t, err, "failure.Validation is expected to return an error")

	assert.Contains(t, err.Error(), failure.ValidationMsg)
}

func TestIsValidation(t *testing.T) {
	err := failure.Validation("some message")
	assert.True(t, failure.IsValidation(err))

	assert.False(t, failure.IsServer(err))
	assert.False(t, failure.IsSystem(err))
	assert.False(t, failure.IsConfig(err))
	assert.False(t, failure.IsNotFound(err))
	assert.False(t, failure.IsIgnore(err))
}

func TestToValidation(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToValidation(e, "user.ID is empty")
	assert.Error(t, err, "failure.Validation is expected to return an error")
	assert.True(t, failure.IsValidation(err))

	expected := fmt.Sprintf("user.ID is empty: api specific msg: %s", failure.ValidationMsg)
	assert.Equal(t, err.Error(), expected)
}

func TestInput(t *testing.T) {
	logMsg := "Invalid User type user.ID is empty"
	lErr := failure.Validation(logMsg)

	inputErr := failure.Input(lErr, "user.ID (%s) is invalid", "xyz")

	assert.Error(t, inputErr, "failure.Input is expected to return an error")

	expected := fmt.Sprintf("Invalid User type user.ID is empty: %s", failure.ValidationMsg)
	assert.Equal(t, expected, inputErr.Error())

	assert.True(t, failure.IsInput(inputErr))
	assert.False(t, failure.IsInput(lErr))

	msg, ok := failure.InputMsg(inputErr)
	require.True(t, ok)
	require.Equal(t, "user.ID (xyz) is invalid", msg)
}

func TestDefer(t *testing.T) {
	err := failure.Defer("some error inside a defer")

	assert.Error(t, err, "failure.Defer is expected to return an error")

	expected := fmt.Sprintf("some error inside a defer: %s", failure.DeferMsg)
	assert.Equal(t, expected, err.Error())

	assert.True(t, failure.IsDefer(err))
	assert.False(t, failure.IsInput(err))

	other := errors.New("some outside err")
	err = failure.ToDefer(other, "something is wrong")

	assert.True(t, failure.IsDefer(err))
}
