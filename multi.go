package failure

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Multi struct {
	Failures  []error
	Formatter MultiFormatFn
}

func (e *Multi) Error() string {
	fn := e.Formatter
	if fn == nil {
		fn = ListFormatFn
	}

	return fn(e.Failures)
}

// ErrorOrNil returns an error interface if this Error represents
// a list of errors, or returns nil if the list of errors is empty. This
// function is useful at the end of accumulation to make sure that the value
// returned represents the existence of errors.
func (e *Multi) ErrorOrNil() error {
	if e == nil {
		return nil
	}
	if len(e.Failures) == 0 {
		return nil
	}

	return e
}

// WrappedErrors returns the list of errors that this Error is wrapping. It is
// an implementation of the errwrap.Wrapper interface so that failure.Multi
// can be used with that library.
//
// This method is not safe to be called concurrently. Unlike accessing the
// Failures field directly, this function also checks if the Multi is nil to
// prevent a null-pointer panic. It satisfies the errwrap.Wrapper interface.
func (e *Multi) WrappedErrors() []error {
	if e == nil {
		return nil
	}
	return e.Failures
}

// Unwrap returns an error from Multi (or nil if there are no errors).
// This error returned will further support Unwrap to get the next error,
// etc. The order will match the order of Errors in the failure.Multi
// at the time of calling.
//
// The resulting error supports errors.As/Is/Unwrap, so you can continue
// to use the stdlib errors package to introspect further.
//
// This will perform a shallow copy of the errors slice. Any errors appended
// to this error after calling Unwrap will not be available until a new
// Unwrap is called on the failure.Multi.
func (e *Multi) Unwrap() error {
	// If we have no errors then we do nothing
	if e == nil || len(e.Failures) == 0 {
		return nil
	}

	// If we have exactly one error, we can just return that directly.
	if len(e.Failures) == 1 {
		return e.Failures[0]
	}

	// Shallow copy the slice
	errs := make([]error, len(e.Failures))
	copy(errs, e.Failures)
	return chain(errs)
}

// chain implements the interfaces necessary for errors.Is/As/Unwrap to
// work in a deterministic way with multierror. A chain tracks a list of
// errors while accounting for the current represented error. This lets
// Is/As be meaningful.
//
// Unwrap returns the next error. In the cleanest form, Unwrap would return
// the wrapped error here but we can't do that if we want to properly
// get access to all the errors. Instead, users are recommended to use
// Is/As to get the correct error type out.
//
// Precondition: []error is non-empty (len > 0)
type chain []error

// Error implements the error interface
func (e chain) Error() string {
	return e[0].Error()
}

// Unwrap implements errors.Unwrap by returning the next error in the
// chain or nil if there are no more errors.
func (e chain) Unwrap() error {
	if len(e) == 1 {
		return nil
	}

	return e[1:]
}

// As implements errors.As by attempting to map to the current value.
func (e chain) As(target interface{}) bool {
	return errors.As(e[0], target)
}

// Is implements errors.Is by comparing the current value directly.
func (e chain) Is(target error) bool {
	return errors.Is(e[0], target)
}

func Append(err error, errs ...error) *Multi {
	switch err := err.(type) {
	case *Multi:
		// Typed nils can be reached here, so initialize if we are nil
		if err == nil {
			err = new(Multi)
		}

		// flat each error
		for _, e := range errs {
			switch e := e.(type) {
			case *Multi:
				if e != nil {
					err.Failures = append(err.Failures, e.Failures...)
				}
			default:
				if e != nil {
					err.Failures = append(err.Failures, e)
				}
			}
		}
		return err
	default:
		newErrs := make([]error, 0, len(errs)+1)
		if err != nil {
			newErrs = append(newErrs, err)
		}
		newErrs = append(newErrs, errs...)
		return Append(&Multi{}, newErrs...)
	}
}

func Multiple(errs []error, opt ...MultiFormatFn) *Multi {
	fn := ListFormatFn
	if len(opt) > 0 && opt[0] != nil {
		fn = opt[0]
	}
	return &Multi{Failures: errs, Formatter: fn}
}

func IsMultiple(e error) bool {
	var t *Multi

	if errors.As(e, &t) {
		return true
	}

	l := errors.Unwrap(e)
	if errors.As(l, &t) {
		return true
	}

	return false
}

func MultiResult(e error) ([]error, bool) {
	err, ok := e.(*Multi)
	if !ok {
		return []error{}, false
	}

	return err.Failures, true
}

// ListFormatFn is a basic formatter that outputs the number of errors
// that occurred along with a bullet point list of the errors.
func ListFormatFn(es []error) string {
	if len(es) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n\n", es[0])
	}

	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n\n",
		len(es), strings.Join(points, "\n\t"))
}

type MultiFormatFn func([]error) string

// Flatten flattens the given error, merging any *Errors together into
// a single *Error.
func Flatten(err error) error {
	// If it isn't an *Error, just return the error as-is
	if _, ok := err.(*Multi); !ok {
		return err
	}

	// Otherwise, make the result and flatten away!
	flatErr := new(Multi)
	flatten(err, flatErr)
	return flatErr
}

func flatten(err error, flatErr *Multi) {
	switch err := err.(type) {
	case *Multi:
		for _, e := range err.Failures {
			flatten(e, flatErr)
		}
	default:
		flatErr.Failures = append(flatErr.Failures, err)
	}
}

// Len implements sort.Interface function for length
func (e Multi) Len() int {
	return len(e.Failures)
}

// Swap implements sort.Interface function for swapping elements
func (e Multi) Swap(i, j int) {
	e.Failures[i], e.Failures[j] = e.Failures[j], e.Failures[i]
}

// Less implements sort.Interface function for determining order
func (e Multi) Less(i, j int) bool {
	return e.Failures[i].Error() < e.Failures[j].Error()
}

// Group is a collection of goroutines which return errors that need to be
// coalesced.
type Group struct {
	mutex sync.Mutex
	err   *Multi
	wg    sync.WaitGroup
}

// Go calls the given function in a new goroutine.
//
// If the function returns an error it is added to the group multierror which
// is returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.mutex.Lock()
			g.err = Append(g.err, err)
			g.mutex.Unlock()
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the Multi
func (g *Group) Wait() *Multi {
	g.wg.Wait()
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.err
}
