package failure_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/rsb/failure"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	msg := "some message"
	err := failure.Server(msg)
	assert.Error(t, err, "failure.Server is expected to return an error")
	assert.True(t, failure.IsServer(err))
	assert.Contains(t, err.Error(), failure.ServerMsg)

	// They are conceptually the same but are not equal
	assert.False(t, failure.IsSystem(err))
	assert.False(t, failure.IsPlatform(err))

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
	assert.False(t, failure.IsPlatform(err))

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

func TestPlatform(t *testing.T) {
	msg := "some message"
	err := failure.Platform(msg)
	assert.Error(t, err, "failure.Platform is expected to return an error")
	assert.True(t, failure.IsPlatform(err))

	assert.False(t, failure.IsServer(err))
	assert.False(t, failure.IsSystem(err))

	assert.Contains(t, err.Error(), failure.PlatformMsg)
}

func TestToPlatform(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToPlatform(e, "api x failed")
	assert.Error(t, err, "failure.ToPlatform is expected to return an error")
	assert.True(t, failure.IsPlatform(err))

	expected := fmt.Sprintf("api x failed: api specific msg: %s", failure.PlatformMsg)
	assert.Equal(t, err.Error(), expected)
}

func TestNotFound(t *testing.T) {
	msg := "some message"
	err := failure.NotFound(msg)
	assert.Error(t, err, "failure.NotFound is expected to return an error")
	assert.True(t, failure.IsNotFound(err))

	assert.False(t, failure.IsServer(err))
	assert.False(t, failure.IsSystem(err))
	assert.False(t, failure.IsPlatform(err))

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
	assert.False(t, failure.IsPlatform(err))
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
	assert.False(t, failure.IsPlatform(err))
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
