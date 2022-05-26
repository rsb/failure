// Package failure implements an opaque error pattern based several of the most
// common  types of errors that occur when developing microservices.
package failure

import (
	"errors"
	"fmt"
)

const (
	SystemMsg             = "system failure"
	ServerMsg             = "server failure"
	NotFoundMsg           = "not found failure"
	NotAuthorizedMsg      = "not authorized failure"
	NotAuthenticatedMsg   = "not authenticated failure"
	ForbiddenMsg          = "access is forbidden"
	ValidationMsg         = "validation failure"
	DeferMsg              = "failure occurred inside defer"
	IgnoreMsg             = "ignore failure"
	ConfigMsg             = "config failure"
	InvalidParamMsg       = "invalid param failure"
	ShutdownMsg           = "system shutdown failure"
	TimeoutMsg            = "timeout failure"
	StartupMsg            = "failure occurred during startup"
	PanicMsg              = "panic"
	BadRequestMsg         = "bad request"
	InvalidAPIFieldsMsg   = "http input fields are not valid"
	MissingFromContextMsg = "resource not in context"
	AlreadyExistsMsg      = "duplicate resource already exists"
	OutOfRangeMsg         = "out of range failure"
	WarnMsg               = "warning"
	NoChangeMsg           = "no change has occurred"
	InvalidStateMsg       = "invalid state"

	systemErr             = err(SystemMsg)
	serverErr             = err(ServerMsg)
	shutdownErr           = err(ShutdownMsg)
	configErr             = err(ConfigMsg)
	notFoundErr           = err(NotFoundMsg)
	notAuthorizedErr      = err(NotAuthorizedMsg)
	notAuthenticatedErr   = err(NotAuthenticatedMsg)
	forbiddenErr          = err(ForbiddenMsg)
	validationErr         = err(ValidationMsg)
	invalidParamErr       = err(InvalidParamMsg)
	deferErr              = err(DeferMsg)
	ignoreErr             = err(IgnoreMsg)
	timeoutErr            = err(TimeoutMsg)
	startupErr            = err(StartupMsg)
	panicErr              = err(PanicMsg)
	badRequestErr         = err(BadRequestMsg)
	invalidAPIFieldsErr   = err(InvalidAPIFieldsMsg)
	missingFromContextErr = err(MissingFromContextMsg)
	alreadyExistsErr      = err(AlreadyExistsMsg)
	outOfRangeErr         = err(OutOfRangeMsg)
	warnErr               = err(WarnMsg)
	noChangeErr           = err(NoChangeMsg)
	invalidStateErr       = err(InvalidStateMsg)
)

type err string

func (e err) Error() string {
	return string(e)
}

// InvalidState is used to signal that the resource is not in a valid state
func InvalidState(format string, a ...interface{}) error {
	return Wrap(invalidStateErr, format, a...)
}

func IsInvalidState(e error) bool {
	return errors.Is(e, invalidStateErr)
}

func ToInvalidState(e error, format string, a ...interface{}) error {
	cause := InvalidState(e.Error())
	return Wrap(cause, format, a...)
}

// NoChange is used to signal that if you expected something to change,
// it has not.
func NoChange(format string, a ...interface{}) error {
	return Wrap(noChangeErr, format, a...)
}

func IsNoChange(e error) bool {
	return errors.Is(e, noChangeErr)
}

func ToNoChange(e error, format string, a ...interface{}) error {
	cause := NoChange(e.Error())
	return Wrap(cause, format, a...)
}

// Warn is used to signal that this error is only a warning. It can be
// used instead of ignore to change the log level of a system
func Warn(format string, a ...interface{}) error {
	return Wrap(warnErr, format, a...)
}

func IsWarn(e error) bool {
	return errors.Is(e, warnErr)
}

func ToWarn(e error, format string, a ...interface{}) error {
	cause := Warn(e.Error())
	return Wrap(cause, format, a...)
}

// OutOfRange is used to signal that the offset of a map is invalid or
// some index for a list is incorrect
func OutOfRange(format string, a ...interface{}) error {
	return Wrap(outOfRangeErr, format, a...)
}

func IsOutOfRange(e error) bool {
	return errors.Is(e, outOfRangeErr)
}

func ToOutOfRange(e error, format string, a ...interface{}) error {
	cause := OutOfRange(e.Error())
	return Wrap(cause, format, a...)
}

// Panic is used in panic recovery blocks or to indicate that you should
// panic if you receive this error
func Panic(format string, a ...interface{}) error {
	return Wrap(panicErr, format, a...)
}

func IsPanic(e error) bool {
	return errors.Is(e, panicErr)
}

func ToPanic(e error, format string, a ...interface{}) error {
	cause := Panic(e.Error())
	return Wrap(cause, format, a...)
}

// MissingFromContext is used to indicate a resource was supposed to be in the
// context but is missing
func MissingFromContext(format string, a ...interface{}) error {
	return Wrap(missingFromContextErr, format, a...)
}

func IsMissingFromContext(e error) bool {
	return errors.Is(e, missingFromContextErr)
}

func ToMissingFromContext(e error, format string, a ...interface{}) error {
	cause := MissingFromContext(e.Error())
	return Wrap(cause, format, a...)
}

// AlreadyExists is used to indicate that the given resource already exists
func AlreadyExists(format string, a ...interface{}) error {
	return Wrap(alreadyExistsErr, format, a...)
}

func IsAlreadyExists(e error) bool {
	return errors.Is(e, alreadyExistsErr)
}

func ToAlreadyExists(e error, format string, a ...interface{}) error {
	cause := AlreadyExists(e.Error())
	return Wrap(cause, format, a...)
}

// Startup is used to signify a failure preventing the system from starting up
func Startup(format string, a ...interface{}) error {
	return Wrap(startupErr, format, a...)
}

func IsStartup(e error) bool {
	return errors.Is(e, startupErr)
}

func ToStartup(e error, format string, a ...interface{}) error {
	cause := Startup(e.Error())
	return Wrap(cause, format, a...)
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
