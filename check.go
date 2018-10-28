package check

import (
	"context"
)

// CheckFunc is the check function only return true or false
type CheckFunc func() bool

// CheckErrFunc returns error if found
type CheckErrFunc func() error

type e struct{}

func wrapCheck(do chan e, ctx context.Context,
	cf CheckFunc, noErrChan chan e) chan e {
	ch := make(chan e)
	go func() {
		<-do
		defer close(ch)

		// not passed
		if !cf() {
			return
		}
		select {
		case <-ctx.Done():
		default:
			noErrChan <- e{}
		}
		<-ctx.Done()
	}()
	return ch
}

func wrapCheckWithError(do chan e, ctx context.Context, cf CheckErrFunc) chan error {
	ch := make(chan error)
	go func() {
		<-do
		defer close(ch)

		select {
		case <-ctx.Done():
			return
		case ch <- cf():
		}
		<-ctx.Done()
	}()
	return ch
}

// Passed returns true if all check points passed, otherwise, returns false
func Passed(ctx context.Context, checkPoints []CheckFunc) bool {
	checkNum := len(checkPoints)
	noErrChan := make(chan e, checkNum)
	eventChan := make(chan e, len(checkPoints))

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	do := make(chan e)
	for _, cf := range checkPoints {
		mergeChan(eventChan, wrapCheck(do, cctx, cf, noErrChan), false)
	}
	close(do)

	passChan := make(chan bool)

	go func() {
		defer close(passChan)
		defer close(eventChan)
		defer close(noErrChan)
		for {
			select {
			case <-noErrChan:
				checkNum--
				if checkNum == 0 {
					close(noErrChan)
					if cctx.Err() == nil {
						passChan <- true
					}
				}
				continue

			case <-eventChan:
			case <-ctx.Done():
				// context cancled
			}
			if cctx.Err() == nil {
				passChan <- false
			}
		}
	}()

	return <-passChan
}

// NoError returns the first error it got, if all passed, returns nil
func NoError(ctx context.Context, checkPoints []CheckErrFunc) error {
	checkNum := len(checkPoints)
	errorChan := make(chan error, len(checkPoints))

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	do := make(chan e)
	for _, cf := range checkPoints {
		mergeChan(errorChan, wrapCheckWithError(do, cctx, cf), true)
	}
	close(do)

	errGot := make(chan error)

	go func() {
		defer close(errGot)
		for {
			select {
			case err := <-errorChan:
				checkNum--
				if checkNum == 0 {
					close(errorChan)
					errGot <- nil
				}
				if err != nil && cctx.Err() == nil {
					errGot <- err
				}
			case <-ctx.Done():
				if cctx.Err() == nil {
					errGot <- ctx.Err()
				}
			}
		}
	}()

	return <-errGot
}

func mergeChan(dest, src interface{}, isError bool) {
	go func() {
		if isError {
			dest.(chan error) <- <-src.(chan error)
		} else {
			dest.(chan e) <- <-src.(chan e)
		}
	}()
}
