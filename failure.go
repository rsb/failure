// Package failure implements an opaque error pattern based several of the most
// common  types of errors that occur when developing microservices.
package failure

import "github.com/pkg/errors"

const (
	SystemMsg     = "system failure"
	PlatformMsg   = "platform failure"
	ServerMsg     = "server failure"
	NotFoundMsg   = "not found failure"
	ValidationMsg = "validation failure"
	IgnoreMsg     = "ignore failure"

	systemErr = err(SystemMsg)
	platErr   = err(PlatformMsg)
	serverErr = err(ServerMsg)

	notFoundErr   = err(NotFoundMsg)
	validationErr = err(ValidationMsg)
	ignoreErr     = err(IgnoreMsg)
)

type err string

func (e err) Error() string {
	return string(e)
}

type inputErr struct {
	public error
	log    error
}

func (e inputErr) Error() string {
	return e.log.Error()
}

func Input(publicErr, internalErr error) error {
	return inputErr{
		public: publicErr,
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

	return i.public.Error(), true
}

// Ignore is used to signify that error should not be acted on, its up
// to the handler to decide to log these errors or not.
func Ignore(format string, a ...interface{}) error {
	return Wrap(ignoreErr, format, a...)
}

func IsIgnore(err error) bool {
	return errors.Cause(err) == ignoreErr
}

// ToIgnore converts `e` into the root cause of ignoreErr, it informs the
// system to ignore this outside error.
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

// ToNotFound converts `e` into the root cause of errNotFound
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

// ToValidation converts `e` into the root cause of errNotFound
func ToValidation(e error, format string, a ...interface{}) error {
	cause := Validation(e.Error())
	return Wrap(cause, format, a...)
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

// ToServer behaves in the same manner as any ToXX function. Making err
// the root cause of type systemErr in our wrap.
func ToServer(e error, format string, a ...interface{}) error {
	cause := Server(e.Error())
	return Wrap(cause, format, a...)
}

// System is has the same meaning as Platform or Server, it can be used instead if you
// don't like how Platform reads in your code
func System(format string, a ...interface{}) error {
	return Wrap(systemErr, format, a...)
}

// IsSystem will return true if the cause is a serverErr
func IsSystem(err error) bool {
	return errors.Cause(err) == systemErr
}

// ToSystem behaves in the same manner as any ToXX function. Making err
// the root cause of type systemErr in our wrap.
func ToSystem(e error, format string, a ...interface{}) error {
	cause := System(e.Error())
	return Wrap(cause, format, a...)
}

// Platform failure is intended to represent low level errors that originate at
// the most concrete part of your architecture. I typically name this layer as
// platform hence the error type. It has the same meaning as Server or System.
// The idea is you choose the name that reads the best for your code and stay
// with that one. It is not recommended mixing these.
func Platform(format string, a ...interface{}) error {
	return Wrap(platErr, format, a...)
}

// IsPlatform will return true if the cause is a platErr
func IsPlatform(err error) bool {
	return errors.Cause(err) == platErr
}

// ToPlatform will convert err to the root cause as a Platform type. This is
// used when the `Platform Layer` of your code any library other than your own,
// and you want that error to be the root cause as type Platform. This will
// preserve the original message, just moving it into the new type.
func ToPlatform(e error, format string, a ...interface{}) error {
	cause := Platform(e.Error())
	return Wrap(cause, format, a...)
}

// Wrap expose errors.Wrapf as our default wrapping style
func Wrap(err error, msg string, a ...interface{}) error {
	return errors.Wrapf(err, msg, a...)
}
