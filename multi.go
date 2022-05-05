package failure

import (
	"fmt"
	"strings"
	"sync"
)

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

func Multiple(errs []error, opt ...MultiFormatFn) error {
	fn := ListFormatFn
	if len(opt) > 0 && opt[0] != nil {
		fn = opt[0]
	}
	return &Multi{Failures: errs, Formatter: fn}
}

func IsMultiple(e error) bool {
	switch e.(type) {
	case *Multi:
		return true
	default:
		return false
	}
}

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
