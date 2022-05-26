package failure_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rsb/failure"
	"github.com/stretchr/testify/assert"
)

func TestInvalidState(t *testing.T) {
	msg := "something is not right"
	err := failure.InvalidState(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), failure.InvalidStateMsg)

	assert.True(t, failure.IsInvalidState(err))
	assert.False(t, failure.IsInvalidState(errors.New("something else")))
}

func TestToInvalidState(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToInvalidState(e, "its not right")
	assert.Error(t, err)
	assert.True(t, failure.IsInvalidState(err))

	expected := "its not right: api specific msg: " + failure.InvalidStateMsg
	assert.Equal(t, err.Error(), expected)
}

func TestNoChange(t *testing.T) {
	msg := "data has not changed"
	err := failure.NoChange(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), failure.NoChangeMsg)

	assert.True(t, failure.IsNoChange(err))
	assert.False(t, failure.IsNoChange(errors.New("something else")))
}

func TestToNoChange(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToNoChange(e, "nothing to do")
	assert.Error(t, err)
	assert.True(t, failure.IsNoChange(err))

	expected := "nothing to do: api specific msg: " + failure.NoChangeMsg
	assert.Equal(t, err.Error(), expected)
}

func TestWarn(t *testing.T) {
	msg := "not really important"
	err := failure.Warn(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), failure.WarnMsg)

	assert.True(t, failure.IsWarn(err))
	assert.False(t, failure.IsWarn(errors.New("something else")))
}

func TestToWarn(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToWarn(e, "should be fixed in the future")
	assert.Error(t, err)
	assert.True(t, failure.IsWarn(err))

	expected := "should be fixed in the future: api specific msg: " + failure.WarnMsg
	assert.Equal(t, err.Error(), expected)
}

func TestOutOfRange(t *testing.T) {
	msg := "invalid offset"
	err := failure.OutOfRange(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), failure.OutOfRangeMsg)

	assert.True(t, failure.IsOutOfRange(err))
	assert.False(t, failure.IsOutOfRange(errors.New("something else")))
}

func TestToOutOfRange(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToOutOfRange(e, "invalid index")
	assert.Error(t, err)
	assert.True(t, failure.IsOutOfRange(err))

	expected := "invalid index: api specific msg: " + failure.OutOfRangeMsg
	assert.Equal(t, err.Error(), expected)
}

func TestMissingFromContext(t *testing.T) {
	msg := "some message"
	err := failure.MissingFromContext(msg)
	assert.Error(t, err, "failure.MissingFromContext is expected to return an error")
	assert.Contains(t, err.Error(), failure.MissingFromContextMsg)

	assert.True(t, failure.IsMissingFromContext(err))
	assert.False(t, failure.IsMissingFromContext(errors.New("something else")))
}

func TestToMissingFromContext(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToMissingFromContext(e, "where is the claim")
	assert.Error(t, err, "failure.ToMissingFromContext is expected to return an error")
	assert.True(t, failure.IsMissingFromContext(err))

	expected := "where is the claim: api specific msg: " + failure.MissingFromContextMsg
	assert.Equal(t, err.Error(), expected)
}

func TestPanic(t *testing.T) {
	msg := "some message"
	err := failure.Panic(msg)
	assert.Error(t, err, "failure.Panic is expected to return an error")
	assert.Contains(t, err.Error(), failure.PanicMsg)

	assert.True(t, failure.IsPanic(err))
	assert.False(t, failure.IsPanic(errors.New("something else")))
}

func TestToPanic(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToPanic(e, "this is not good")
	assert.Error(t, err, "failure.ToPanic is expected to return an error")
	assert.True(t, failure.IsPanic(err))

	expected := "this is not good: api specific msg: " + failure.PanicMsg
	assert.Equal(t, err.Error(), expected)
}

func TestAlreadyExists(t *testing.T) {
	msg := "some message"
	err := failure.AlreadyExists(msg)
	assert.Error(t, err, "failure.Already is expected to return an error")
	assert.Contains(t, err.Error(), failure.AlreadyExistsMsg)

	assert.True(t, failure.IsAlreadyExists(err))
	assert.False(t, failure.IsAlreadyExists(errors.New("something else")))
}

func TestToAlreadyExists(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToAlreadyExists(e, "this is duplicated")
	assert.Error(t, err, "failure.ToAlreadyExists is expected to return an error")
	assert.True(t, failure.IsAlreadyExists(err))

	expected := "this is duplicated: api specific msg: " + failure.AlreadyExistsMsg
	assert.Equal(t, err.Error(), expected)
}

func TestStartup(t *testing.T) {
	msg := "some message"
	err := failure.Startup(msg)
	assert.Error(t, err, "failure.Startup is expected to return an error")
	assert.Contains(t, err.Error(), failure.StartupMsg)

	assert.True(t, failure.IsStartup(err))
	assert.False(t, failure.IsTimeout(err))

}

func TestToStartup(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToStartup(e, "initialized wrong")
	assert.Error(t, err, "failure.ToStartup is expected to return an error")
	assert.True(t, failure.IsStartup(err))

	expected := "initialized wrong: api specific msg: " + failure.StartupMsg
	assert.Equal(t, err.Error(), expected)
}

func TestTimeout(t *testing.T) {
	msg := "some message"
	err := failure.Timeout(msg)
	assert.Error(t, err, "failure.Timeout is expected to return an error")
	assert.Contains(t, err.Error(), failure.TimeoutMsg)

	assert.True(t, failure.IsTimeout(err))
	assert.False(t, failure.IsNotAuthorized(err))

}

func TestToTimeout(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToTimeout(e, "this took way too long")
	assert.Error(t, err, "failure.ToTimeout is expected to return an error")
	assert.True(t, failure.IsTimeout(err))

	expected := "this took way too long: api specific msg: " + failure.TimeoutMsg
	assert.Equal(t, err.Error(), expected)
}

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

func TestDefer(t *testing.T) {
	err := failure.Defer("some error inside a defer")

	assert.Error(t, err, "failure.Defer is expected to return an error")

	expected := fmt.Sprintf("some error inside a defer: %s", failure.DeferMsg)
	assert.Equal(t, expected, err.Error())

	assert.True(t, failure.IsDefer(err))
	assert.False(t, failure.IsNotAuthorized(err))

	other := errors.New("some outside err")
	err = failure.ToDefer(other, "something is wrong")

	assert.True(t, failure.IsDefer(err))
}
