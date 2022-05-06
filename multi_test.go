package failure_test

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/rsb/failure"
)

func TestMulti_Impl(t *testing.T) {
	var _ error = new(failure.Multi)
}

func TestErrorError_custom(t *testing.T) {
	errs := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	fn := func(es []error) string {
		return "foo"
	}

	multi := &failure.Multi{Failures: errs, Formatter: fn}
	assert.Equal(t, multi.Error(), "foo")
	assert.True(t, failure.IsMultiple(multi))
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
	assert.True(t, failure.IsMultiple(multi))

	require.True(t, reflect.DeepEqual(multi.Failures, multi.WrappedErrors()))

	multi = nil
	require.Nil(t, multi.WrappedErrors())
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

func TestErrorIs(t *testing.T) {
	errBar := errors.New("bar")

	t.Run("with errBar", func(t *testing.T) {
		err := &failure.Multi{Failures: []error{
			errors.New("foo"),
			errBar,
			errors.New("baz"),
		}}
		require.True(t, errors.Is(err, errBar))

	})

	t.Run("with errBar wrapped by fmt.Errorf", func(t *testing.T) {
		err := &failure.Multi{Failures: []error{
			errors.New("foo"),
			fmt.Errorf("errorf: %w", errBar),
			errors.New("baz"),
		}}

		require.True(t, errors.Is(err, errBar))
	})

	t.Run("without errBar", func(t *testing.T) {
		err := &failure.Multi{Failures: []error{
			errors.New("foo"),
			errors.New("baz"),
		}}

		require.False(t, errors.Is(err, errBar))
	})
}

func TestErrorAs(t *testing.T) {
	match := &nestedError{}

	t.Run("with the value", func(t *testing.T) {
		err := &failure.Multi{Failures: []error{
			errors.New("foo"),
			match,
			errors.New("baz"),
		}}

		var target *nestedError
		if !errors.As(err, &target) {
			t.Fatal("should be true")
		}
		if target == nil {
			t.Fatal("target should not be nil")
		}
	})

	t.Run("with the value wrapped by fmt.Errorf", func(t *testing.T) {
		err := &failure.Multi{Failures: []error{
			errors.New("foo"),
			fmt.Errorf("errorf: %w", match),
			errors.New("baz"),
		}}

		var target *nestedError
		require.True(t, errors.As(err, &target))
		require.NotNil(t, target)
	})

	t.Run("without the value", func(t *testing.T) {
		err := &failure.Multi{Failures: []error{
			errors.New("foo"),
			errors.New("baz"),
		}}

		var target *nestedError
		require.False(t, errors.As(err, &target))
		require.Nil(t, target)
	})
}

// nestedError implements error and is used for tests.
type nestedError struct{}

func (*nestedError) Error() string { return "" }

func Test_MultiFlatten(t *testing.T) {
	original := &failure.Multi{
		Failures: []error{
			errors.New("one"),
			&failure.Multi{
				Failures: []error{
					errors.New("two"),
					&failure.Multi{
						Failures: []error{
							errors.New("three"),
						},
					},
				},
			},
		},
	}

	expected := `3 errors occurred:
	* one
	* two
	* three

`
	actual := fmt.Sprintf("%s", failure.Flatten(original))

	assert.Equal(t, expected, actual)
}

func Test_MultiFlattenNonError(t *testing.T) {
	err := errors.New("foo")
	actual := failure.Flatten(err)
	assert.True(t, reflect.DeepEqual(actual, err))
}

func Test_Multi_SortSingle(t *testing.T) {
	errFoo := errors.New("foo")

	expected := []error{
		errFoo,
	}

	err := &failure.Multi{
		Failures: []error{
			errFoo,
		},
	}

	sort.Sort(err)
	assert.True(t, reflect.DeepEqual(err.Failures, expected))
}

func Test_Multi_SortMultiple(t *testing.T) {
	errBar := errors.New("bar")
	errBaz := errors.New("baz")
	errFoo := errors.New("foo")

	expected := []error{
		errBar,
		errBaz,
		errFoo,
	}

	err := &failure.Multi{
		Failures: []error{
			errFoo,
			errBar,
			errBaz,
		},
	}

	sort.Sort(err)
	assert.True(t, reflect.DeepEqual(err.Failures, expected))
}

func Test_Multi_Group(t *testing.T) {
	err1 := errors.New("group_test: 1")
	err2 := errors.New("group_test: 2")

	cases := []struct {
		errs      []error
		nilResult bool
	}{
		{errs: []error{}, nilResult: true},
		{errs: []error{nil}, nilResult: true},
		{errs: []error{err1}},
		{errs: []error{err1, nil}},
		{errs: []error{err1, nil, err2}},
	}

	for _, tc := range cases {
		var g failure.Group

		for _, err := range tc.errs {
			err := err
			g.Go(func() error { return err })

		}

		gErr := g.Wait()
		if gErr != nil {
			for i := range tc.errs {
				if tc.errs[i] != nil && !strings.Contains(gErr.Error(), tc.errs[i].Error()) {
					t.Fatalf("expected error to contain %q, actual: %v", tc.errs[i].Error(), gErr)
				}
			}
		} else if !tc.nilResult {
			t.Fatalf("Group.Wait() should not have returned nil for errs: %v", tc.errs)
		}
	}
}

func TestMultiResult(t *testing.T) {
	list := []error{
		failure.Timeout("some timeout"),
		failure.System("some other error"),
		failure.Ignore("some ignore"),
	}

	err := failure.Multiple(list)
	require.Error(t, err)

	result, ok := failure.MultiResult(err)
	require.True(t, ok)
	require.Equal(t, list, result)

	e := errors.New("some other thing")
	result, ok = failure.MultiResult(e)
	require.False(t, ok)
	require.Empty(t, result)
}
