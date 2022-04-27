package failure_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"

	"github.com/rsb/failure"
	"github.com/stretchr/testify/assert"
)

func TestIsAnyAuthFailure(t *testing.T) {
	err := failure.NotFound("something not found")
	err1 := failure.Forbidden("some message")
	err2 := failure.NotAuthorized("you are not authorized")
	err3 := failure.NotAuthenticated("you are not authenticated")

	assert.False(t, failure.IsAnyAuthFailure(err))
	assert.True(t, failure.IsAnyAuthFailure(err1))
	assert.True(t, failure.IsAnyAuthFailure(err2))
	assert.True(t, failure.IsAnyAuthFailure(err3))

}

func TestForbidden(t *testing.T) {
	msg := "some message"
	err := failure.Forbidden(msg)
	assert.Error(t, err, "failure.Forbidden is expected to return an error")
	assert.Contains(t, err.Error(), failure.ForbiddenMsg)

	assert.True(t, failure.IsForbidden(err))
	assert.False(t, failure.IsNotAuthorized(err))

}

func TestToForbidden(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToForbidden(e, "not wanted here")
	assert.Error(t, err, "failure.ToForbidden is expected to return an error")
	assert.True(t, failure.IsForbidden(err))

	expected := "not wanted here: api specific msg: " + failure.ForbiddenMsg
	assert.Equal(t, err.Error(), expected)
}

func TestNotAuthenticated(t *testing.T) {
	msg := "some message"
	err := failure.NotAuthenticated(msg)
	assert.Error(t, err, "failure.NotAuthenticated is expected to return an error")
	assert.Contains(t, err.Error(), failure.NotAuthenticatedMsg)

	assert.True(t, failure.IsNotAuthenticated(err))
	assert.False(t, failure.IsNotAuthorized(err))

}

func TestToNotAuthenticated(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToNotAuthenticated(e, "access denied")
	assert.Error(t, err, "failure.ToNotAuthenticated is expected to return an error")
	assert.True(t, failure.IsNotAuthenticated(err))

	expected := "access denied: api specific msg: " + failure.NotAuthenticatedMsg
	assert.Equal(t, err.Error(), expected)
}

func TestNotAuthorized(t *testing.T) {
	msg := "some message"
	err := failure.NotAuthorized(msg)
	assert.Error(t, err, "failure.NotAuthorized is expected to return an error")
	assert.Contains(t, err.Error(), failure.NotAuthorizedMsg)

	assert.True(t, failure.IsNotAuthorized(err))
	assert.False(t, failure.IsSystem(err))

}

func TestToNotAuthorized(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToNotAuthorized(e, "user not allowed")
	assert.Error(t, err, "failure.ToNotAuthorized is expected to return an error")
	assert.True(t, failure.IsNotAuthorized(err))

	expected := "user not allowed: api specific msg: " + failure.NotAuthorizedMsg
	assert.Equal(t, err.Error(), expected)
}

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
	fields := map[string]string{
		"field1": "invalid option 1",
		"field2": "invalid option 2",
	}
	inputErr := failure.InvalidInput(fields, "data given has the following errors")

	assert.Error(t, inputErr, "failure.Input is expected to return an error")

	expected := "invalid input: data given has the following errors: field1: invalid option 1,field2: invalid option 2"
	assert.Equal(t, expected, inputErr.Error())

	assert.True(t, failure.IsInvalidInput(inputErr))

	otherErr := errors.New("some other error")
	assert.False(t, failure.IsInvalidInput(otherErr))

	result, ok := failure.InputFields(inputErr)
	require.True(t, ok)
	assert.Equal(t, fields, result)

	out, ok := failure.InvalidInputMsg(inputErr)
	require.True(t, ok)
	assert.Equal(t, "data given has the following errors", out)

}

func TestDefer(t *testing.T) {
	err := failure.Defer("some error inside a defer")

	assert.Error(t, err, "failure.Defer is expected to return an error")

	expected := fmt.Sprintf("some error inside a defer: %s", failure.DeferMsg)
	assert.Equal(t, expected, err.Error())

	assert.True(t, failure.IsDefer(err))
	assert.False(t, failure.IsInvalidInput(err))

	other := errors.New("some outside err")
	err = failure.ToDefer(other, "something is wrong")

	assert.True(t, failure.IsDefer(err))
}
