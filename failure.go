// Package failure implements an opaque error pattern based several of the most
// common  types of errors that occur when developing microservices.
package failure

import (
	"errors"
	"fmt"
)

const (
	SystemMsg           = "system failure"
	ServerMsg           = "server failure"
	NotFoundMsg         = "not found failure"
	NotAuthorizedMsg    = "not authorized failure"
	NotAuthenticatedMsg = "not authenticated failure"
	ForbiddenMsg        = "access is forbidden"
	ValidationMsg       = "validation failure"
	DeferMsg            = "failure occurred inside defer"
	IgnoreMsg           = "ignore failure"
	ConfigMsg           = "config failure"
	InvalidParamMsg     = "invalid param failure"
	ShutdownMsg         = "system shutdown failure"
	BadRequestMsg       = "bad request"
	TimeoutMsg          = "timeout failure"
	NoMoreRetriesMsg    = "retry limit reached"

	systemErr           = err(SystemMsg)
	serverErr           = err(ServerMsg)
	shutdownErr         = err(ShutdownMsg)
	configErr           = err(ConfigMsg)
	notFoundErr         = err(NotFoundMsg)
	notAuthorizedErr    = err(NotAuthorizedMsg)
	notAuthenticatedErr = err(NotAuthenticatedMsg)
	forbiddenErr        = err(ForbiddenMsg)
	validationErr       = err(ValidationMsg)
	invalidParamErr     = err(InvalidParamMsg)
	deferErr            = err(DeferMsg)
	ignoreErr           = err(IgnoreMsg)
	badRequestErr       = err(BadRequestMsg)
	timeoutErr          = err(TimeoutMsg)
	noMoreRetriesErr    = err(NoMoreRetriesMsg)
)

type err string

func (e err) Error() string {
	return string(e)
}

type ignoreRetryErr struct {
	err error
}

func (e ignoreRetryErr) Error() string {
	return e.err.Error()
}

func (e ignoreRetryErr) Unwrap() error {
	return e.err
}

// Timeout is used to signify that error because something was taking
// too long
func Timeout(format string, a ...interface{}) error {
	return Wrap(timeoutErr, format, a...)
}

func IsTimeout(e error) bool {
	return errors.Is(e, timeoutErr)
}

func ToTimeout(e error, format string, a ...interface{}) error {
	cause := Timeout(e.Error())
	return Wrap(cause, format, a...)
}

// NoMoreRetries is used to signal all available retries have been used up
func NoMoreRetries(format string, a ...interface{}) error {
	return Wrap(noMoreRetriesErr, format, a...)
}

func IsNoMoreRetries(e error) bool {
	return errors.Is(e, noMoreRetriesErr)
}

func ToNoMoreRetries(e error, format string, a ...interface{}) error {
	cause := NoMoreRetries(e.Error())
	return Wrap(cause, format, a...)
}

// IgnoreRetry is used to signal that this operation should not be
// retried. Useful in http client middleware that provides operations
// that are mixed into retry or backoff functions
func IgnoreRetry(format string, a ...interface{}) error {
	return ignoreRetryErr{
		err: errors.New(fmt.Sprintf(format, a...)),
	}
}

func IsIgnoreRetry(e error) bool {
	var t ignoreRetryErr
	if found := errors.As(e, &t); !found {
		return false
	}

	return true
}

func ToIgnoreRetry(e error, format string, a ...interface{}) error {
	cause := IgnoreRetry(e.Error())
	return Wrap(cause, format, a...)
}

func WrapIgnoreRetry(e error) error {
	return ignoreRetryErr{
		err: e,
	}
}

func UnwrapIgnoreRetry(e error) error {
	var t ignoreRetryErr

	if found := errors.As(e, &t); !found {
		return e
	}

	return t.err
}

// Config is used to signify that error occurred when processing the
// application configuration
func Config(format string, a ...interface{}) error {
	return Wrap(configErr, format, a...)
}

func IsConfig(e error) bool {
	return errors.Is(e, configErr)
}

func ToConfig(e error, format string, a ...interface{}) error {
	cause := Config(e.Error())
	return Wrap(cause, format, a...)
}

// InvalidParam is to indicate that the param of a function or any
// parameter in general is invalid
func InvalidParam(format string, a ...interface{}) error {
	return Wrap(invalidParamErr, format, a...)
}

func IsInvalidParam(e error) bool {
	return errors.Is(e, invalidParamErr)
}

func ToInvalidParam(e error, format string, a ...interface{}) error {
	cause := InvalidParam(e.Error())
	return Wrap(cause, format, a...)
}

// Ignore is used to signify that error should not be acted on, it's up
// to the handler to decide to log these errors or not.
func Ignore(format string, a ...interface{}) error {
	return Wrap(ignoreErr, format, a...)
}

func IsIgnore(e error) bool {
	return errors.Is(e, ignoreErr)
}

// ToIgnore converts `e` into the root cause of ignoreErr, it informs the
// system to ignore error. Used typically to log results and do not act on
// the error itself.
func ToIgnore(e error, format string, a ...interface{}) error {
	cause := Ignore(e.Error())
	return Wrap(cause, format, a...)
}

// NotFound is used to signify that whatever resource you were looking for
// does not exist and that fact it does not exist is an error.
func NotFound(format string, a ...interface{}) error {
	return Wrap(notFoundErr, format, a...)
}

func IsNotFound(e error) bool {
	return errors.Is(e, notFoundErr)
}

func ToNotFound(e error, format string, a ...interface{}) error {
	cause := NotFound(e.Error())
	return Wrap(cause, format, a...)
}

// NotAuthorized is used to signify that a resource does not have sufficient
// access to perform a given task
func NotAuthorized(format string, a ...interface{}) error {
	return Wrap(notAuthorizedErr, format, a...)
}

func IsNotAuthorized(e error) bool {
	return errors.Is(e, notAuthorizedErr)
}

func ToNotAuthorized(e error, format string, a ...interface{}) error {
	cause := NotAuthorized(e.Error())
	return Wrap(cause, format, a...)
}

// NotAuthenticated is used to signify that a resource's identity verification
// failed. They are not who they claim to be
func NotAuthenticated(format string, a ...interface{}) error {
	return Wrap(notAuthenticatedErr, format, a...)
}

func IsNotAuthenticated(e error) bool {
	return errors.Is(e, notAuthenticatedErr)
}

func ToNotAuthenticated(e error, format string, a ...interface{}) error {
	cause := NotAuthenticated(e.Error())
	return Wrap(cause, format, a...)
}

// Forbidden is used to signify either not authenticated or
// not authorized
func Forbidden(format string, a ...interface{}) error {
	return Wrap(forbiddenErr, format, a...)
}

func IsForbidden(e error) bool {
	return errors.Is(e, forbiddenErr)
}

func ToForbidden(e error, format string, a ...interface{}) error {
	cause := Forbidden(e.Error())
	return Wrap(cause, format, a...)
}

// IsAnyAuthFailure can be used to determine if any of the following we used:
// NotAuthenticated, NotAuthorized, Forbidden
func IsAnyAuthFailure(e error) bool {
	return IsNotAuthenticated(e) ||
		IsNotAuthorized(e) ||
		IsForbidden(e)
}

// Validation is used to signify that a validation rule as been violated
func Validation(format string, a ...interface{}) error {
	return Wrap(validationErr, format, a...)
}

func IsValidation(e error) bool {
	return errors.Is(e, validationErr)
}

func ToValidation(e error, format string, a ...interface{}) error {
	cause := Validation(e.Error())
	return Wrap(cause, format, a...)
}

// Defer is used to signify errors that originate inside a defer function
func Defer(format string, a ...interface{}) error {
	return Wrap(deferErr, format, a...)
}

func IsDefer(e error) bool {
	return errors.Is(e, deferErr)
}

func ToDefer(e error, format string, a ...interface{}) error {
	cause := Defer(e.Error())
	return Wrap(cause, format, a...)
}

// Shutdown is used to signal that the app should shut down.
func Shutdown(format string, a ...interface{}) error {
	return Wrap(shutdownErr, format, a...)
}

func ToShutdown(e error, format string, a ...interface{}) error {
	cause := Shutdown(e.Error())
	return Wrap(cause, format, a...)
}

func IsShutdown(e error) bool {
	return errors.Is(e, shutdownErr)
}

// BadRequest is used to signal that the app should shut down.
func BadRequest(format string, a ...interface{}) error {
	return Wrap(badRequestErr, format, a...)
}

func ToBadRequest(e error, format string, a ...interface{}) error {
	cause := BadRequest(e.Error())
	return Wrap(cause, format, a...)
}

func IsBadRequest(e error) bool {
	return errors.Is(e, badRequestErr)
}

// Server has the same meaning as Platform or System, it can be used instead if you
// don't like how Platform or System reads in your code.
func Server(format string, a ...interface{}) error {
	return Wrap(serverErr, format, a...)
}

// IsServer will return true if the cause is a serverErr
func IsServer(err error) bool {
	return errors.Is(err, serverErr)
}

func ToServer(e error, format string, a ...interface{}) error {
	cause := Server(e.Error())
	return Wrap(cause, format, a...)
}

// System is has the same meaning as Platform or Server, it can be used instead if you
// don't like how Platform reads in your code
func System(format string, a ...interface{}) error {
	return Wrap(systemErr, format, a...)
}

func IsSystem(err error) bool {
	return errors.Is(err, systemErr)
}

func ToSystem(e error, format string, a ...interface{}) error {
	cause := System(e.Error())
	return Wrap(cause, format, a...)
}

// Wrap expose errors.Wrapf as our default wrapping style
func Wrap(err error, msg string, a ...interface{}) error {
	msg = fmt.Sprintf(msg, a...)
	return fmt.Errorf("%s: %w", msg, err)
}
