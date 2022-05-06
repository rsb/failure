package failure_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/rsb/failure"
)

func TestMulti_Impl(t *testing.T) {
	var _ error = new(failure.Multi)
}

func TestErrorError_custom(t *testing.T) {
	errors := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	fn := func(es []error) string {
		return "foo"
	}

	multi := &failure.Multi{Failures: errors, Formatter: fn}
	assert.Equal(t, multi.Error(), "foo")
}

func TestErrorError_default(t *testing.T) {
	expected := `2 errors occurred:
	* foo
	* bar

`

	list := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	multi := failure.Multiple(list)
	assert.Equal(t, expected, multi.Error())
}

func TestErrorErrorOrNil(t *testing.T) {
	err := new(failure.Multi)
	require.NoError(t, err.ErrorOrNil())

	err.Failures = []error{errors.New("foo")}

	v := err.ErrorOrNil()
	require.Error(t, v)
	require.Equal(t, v, err)
}

func TestErrorWrappedErrors(t *testing.T) {
	list := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	multi := failure.Multiple(list)
	if !reflect.DeepEqual(multi.Failures, multi.WrappedErrors()) {
		t.Fatalf("bad: %s", multi.WrappedErrors())
	}

	multi = nil
	if err := multi.WrappedErrors(); err != nil {
		t.Fatalf("bad: %#v", multi)
	}
}

func TestErrorUnwrap(t *testing.T) {
	t.Run("with errors", func(t *testing.T) {
		err := &failure.Multi{Failures: []error{
			errors.New("foo"),
			errors.New("bar"),
			errors.New("baz"),
		}}

		var current error = err
		for i := 0; i < len(err.Failures); i++ {
			current = errors.Unwrap(current)
			require.True(t, errors.Is(current, err.Failures[i]))
		}

		require.Nil(t, errors.Unwrap(current))
	})

	t.Run("with no errors", func(t *testing.T) {
		err := &failure.Multi{Failures: nil}
		require.Nil(t, errors.Unwrap(err))
	})

	t.Run("with nil", func(t *testing.T) {
		var err *failure.Multi
		require.Nil(t, errors.Unwrap(err))
	})
}
