package failure_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/rsb/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRest_Error(t *testing.T) {
	r := &failure.RestAPI{
		StatusCode: http.StatusForbidden,
		Msg:        "not allowed",
		Err:        errors.New("some api error"),
	}

	err, ok := failure.RestError(r)
	require.True(t, ok)
	assert.Equal(t, r.Err, err)

	err, ok = failure.RestError(errors.New("some other error"))
	require.False(t, ok)
	assert.Nil(t, err)
}

func TestRest_StatusCode(t *testing.T) {
	r := &failure.RestAPI{
		StatusCode: http.StatusForbidden,
		Msg:        "not allowed",
		Err:        errors.New("some api error"),
	}

	code, ok := failure.RestStatusCode(r)
	require.True(t, ok)
	assert.Equal(t, r.StatusCode, code)

	code, ok = failure.RestStatusCode(errors.New("some other error"))
	require.False(t, ok)
	assert.Equal(t, 0, code)
}

func TestRest_Message(t *testing.T) {
	r := &failure.RestAPI{
		StatusCode: http.StatusForbidden,
		Msg:        "not allowed",
		Err:        errors.New("some api error"),
	}

	msg, ok := failure.RestMessage(r)
	require.True(t, ok)
	assert.Equal(t, r.Msg, msg)

	msg, ok = failure.RestMessage(errors.New("some other error"))
	require.False(t, ok)
	assert.Empty(t, msg)
}

func TestRest_NewInvalidFields(t *testing.T) {
	fields := map[string]string{
		"field-a": "error a",
		"field-b": "error b",
	}
	err := failure.NewInvalidFields(fields, "invalid inputs given")
	require.Error(t, err)

	e := func() error {
		return err
	}()

	out, ok := failure.GetInvalidFields(e)
	require.True(t, ok)
	require.Equal(t, fields, out)

	notRest := errors.New("some error")

	out, ok = failure.GetInvalidFields(notRest)
	require.False(t, ok)
	assert.Nil(t, out)
}

func TestRestAPIFields_Error(t *testing.T) {
	fields := map[string]string{
		"field-a": "error a",
		"field-b": "error b",
	}
	err := failure.InvalidFields(fields, "invalid inputs given")
	require.Error(t, err)

	assert.True(t, failure.IsRestAPI(err))
	assert.True(t, failure.IsInvalidFields(err))

	other := errors.New("foo")
	assert.False(t, failure.IsInvalidFields(other))

	assert.False(t, failure.IsRestAPI(errors.New("foo")))
}

func TestRestAPI_RestError(t *testing.T) {
	err := failure.NewBadRequest("some bad request (%s)", "foo")

	e, ok := failure.RestError(err)
	require.True(t, ok)
	require.Equal(t, err.Err, e)
}

func TestRestAPI_Error(t *testing.T) {
	err := failure.BadRequest("some bad request (%s)", "foo")
	assert.True(t, failure.IsBadRequest(err))
	assert.True(t, failure.IsRestAPI(err))
	code, ok := failure.RestStatusCode(err)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, code)

	returnErr := func() error {
		return err
	}

	e := returnErr()

	assert.True(t, failure.IsBadRequest(e))
	assert.True(t, failure.IsRestAPI(e))

	code, ok = failure.RestStatusCode(e)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, code)

	returnWrap := func() error {
		tmp := failure.Wrap(err, "some wrap")
		return failure.Wrap(tmp, "other message")
	}

	w := returnWrap()
	assert.True(t, failure.IsBadRequest(w))
	assert.True(t, failure.IsRestAPI(w))

	code, ok = failure.RestStatusCode(w)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, code)

	api := errors.New("some outside failure")
	e = failure.ToBadRequest(api, "This field is invalid")
	assert.True(t, failure.IsBadRequest(e))
	assert.True(t, failure.IsRestAPI(e))

	code, ok = failure.RestStatusCode(e)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, code)

}
func TestBadRequest(t *testing.T) {
	msg := "some message"
	err := failure.BadRequest(msg)
	assert.Error(t, err, "failure.BadRequest is expected to return an error")
	assert.True(t, failure.IsBadRequest(err))
	assert.Contains(t, err.Error(), failure.BadRequestMsg)

	assert.False(t, failure.IsSystem(err))
	assert.False(t, failure.IsBadRequest(errors.New("foo")))
}

func TestToBadRequest(t *testing.T) {
	msg := "api specific msg"
	e := errors.New(msg)

	err := failure.ToBadRequest(e, "user messed up")
	assert.Error(t, err, "failure.ToBadRequest is expected to return an error")
	assert.True(t, failure.IsBadRequest(err))

	expected := "api specific msg"
	assert.Equal(t, expected, err.Error())
}
