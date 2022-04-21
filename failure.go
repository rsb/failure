// Package failure implements an opaque error pattern based several of the most
// common  types of errors that occur when developing microservices.
package failure

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	SystemMsg       = "system failure"
	ServerMsg       = "server failure"
	NotFoundMsg     = "not found failure"
	ValidationMsg   = "validation failure"
	DeferMsg        = "failure occurred inside defer"
	IgnoreMsg       = "ignore failure"
	ConfigMsg       = "config failure"
	InvalidParamMsg = "invalid param failure"
	ShutdownMsg     = "system shutdown failure"

	systemErr       = err(SystemMsg)
	serverErr       = err(ServerMsg)
	shutdownErr     = err(ShutdownMsg)
	configErr       = err(ConfigMsg)
	notFoundErr     = err(NotFoundMsg)
	validationErr   = err(ValidationMsg)
	invalidParamErr = err(InvalidParamMsg)
	deferErr        = err(DeferMsg)
	ignoreErr       = err(IgnoreMsg)
)

type err string

func (e err) Error() string {
	return string(e)
}

type inputErr struct {
	public string
	log    error
}

func (e inputErr) Error() string {
	return e.log.Error()
}

func Input(internalErr error, format string, a ...interface{}) error {
	return inputErr{
		public: fmt.Sprintf(format, a...),
		log:    internalErr,
	}
}

func IsInput(e error) bool {
	root := errors.Cause(e)
	if _, ok := root.(inputErr); !ok {
		return false
	}

	return true
}

func InputMsg(e error) (string, bool) {
	root := errors.Cause(e)
	i, ok := root.(inputErr)
	if !ok {
		return "", false
	}

	return i.public, true
}

// Config is used to signify that error occurred when processing the
// application configuration
func Config(format string, a ...interface{}) error {
	return Wrap(configErr, format, a...)
}

func IsConfig(err error) bool {
	return errors.Cause(err) == configErr
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

func IsInvalidParam(err error) bool {
	return errors.Cause(err) == invalidParamErr
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

func IsIgnore(err error) bool {
	return errors.Cause(err) == ignoreErr
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

func IsNotFound(err error) bool {
	return errors.Cause(err) == notFoundErr
}

func ToNotFound(e error, format string, a ...interface{}) error {
	cause := NotFound(e.Error())
	return Wrap(cause, format, a...)
}

// Validation is used to signify that a validation rule as been violated
func Validation(format string, a ...interface{}) error {
	return Wrap(validationErr, format, a...)
}

func IsValidation(err error) bool {
	return errors.Cause(err) == validationErr
}

func ToValidation(e error, format string, a ...interface{}) error {
	cause := Validation(e.Error())
	return Wrap(cause, format, a...)
}

// Defer is used to signify errors that originate inside a defer function
func Defer(format string, a ...interface{}) error {
	return Wrap(deferErr, format, a...)
}

func IsDefer(err error) bool {
	return errors.Cause(err) == deferErr
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
	return errors.Cause(e) == shutdownErr
}

// Server has the same meaning as Platform or System, it can be used instead if you
// don't like how Platform or System reads in your code.
func Server(format string, a ...interface{}) error {
	return Wrap(serverErr, format, a...)
}

// IsServer will return true if the cause is a serverErr
func IsServer(err error) bool {
	return errors.Cause(err) == serverErr
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
	return errors.Cause(err) == systemErr
}

func ToSystem(e error, format string, a ...interface{}) error {
	cause := System(e.Error())
	return Wrap(cause, format, a...)
}

// Wrap expose errors.Wrapf as our default wrapping style
func Wrap(err error, msg string, a ...interface{}) error {
	return errors.Wrapf(err, msg, a...)
}
