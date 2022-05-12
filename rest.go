package failure

import (
	"errors"
	"fmt"
	"net/http"
)

type RestAPI struct {
	StatusCode int
	Msg        string
	Fields     map[string]string
	Err        error
}

func (r *RestAPI) Error() string {
	return r.Err.Error()
}

func NewInvalidFields(f map[string]string, msg string, a ...interface{}) *RestAPI {
	r := RestAPI{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        fmt.Sprintf(msg, a...),
		Fields:     f,
		Err:        invalidAPIFieldsErr,
	}

	return &r
}

func InvalidFields(f map[string]string, msg string, a ...interface{}) error {
	return NewInvalidFields(f, msg, a...)
}

func GetInvalidFields(e error) (map[string]string, bool) {
	var r *RestAPI
	if !errors.As(e, &r) {
		return nil, false
	}

	return r.Fields, true
}

func IsInvalidFields(e error) bool {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.StatusCode == http.StatusUnprocessableEntity
	}

	return false
}

func NewBadRequest(msg string, a ...interface{}) *RestAPI {
	r := RestAPI{
		StatusCode: http.StatusBadRequest,
		Msg:        fmt.Sprintf(msg, a...),
		Err:        badRequestErr,
	}
	return &r
}

func BadRequest(msg string, a ...interface{}) error {
	return NewBadRequest(msg, a...)
}

func ToBadRequest(e error, msg string, a ...interface{}) error {
	r := RestAPI{
		StatusCode: http.StatusBadRequest,
		Msg:        fmt.Sprintf(msg, a...),
		Err:        e,
	}
	return &r
}

func IsBadRequest(e error) bool {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.StatusCode == http.StatusBadRequest
	}

	return false
}

func RestStatusCode(e error) (int, bool) {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.StatusCode, true
	}

	return 0, false
}

func RestMessage(e error) (string, bool) {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.Msg, true
	}

	return "", false
}

func RestError(e error) (error, bool) {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.Err, true
	}

	return nil, false
}

func IsRestAPI(e error) bool {
	var r *RestAPI

	if errors.As(e, &r) {
		return true
	}

	return false
}
