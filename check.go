package check

import (
	"context"
	"sync"
)

var wg sync.WaitGroup

// CheckFunc is the check function only return true or false
type CheckFunc func() bool

// CheckErrFunc returns error if found
type CheckErrFunc func() error

type e struct{}

func wrapCheck(do chan e, ctx context.Context, cf CheckFunc, noErrChan chan e) chan e {
	ch := make(chan e)
	go func() {
		c := make(chan bool, 1)
		<-do

		select {
		case <-ctx.Done():
			close(ch) // not pass
		case c <- cf():
			close(c)
			if <-c {
				noErrChan <- e{}
			} else {
				// not pass
				close(ch)
			}
		}
	}()
	return ch
}

func wrapCheckWithError(do chan e, ctx context.Context, cf CheckErrFunc) chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		<-do

		select {
		case <-ctx.Done():
		case ch <- cf():
		}
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
			case <-eventChan:
				checkNum--
				if cctx.Err() == nil {
					passChan <- false
				}

			case <-ctx.Done():
				if cctx.Err() == nil {
					passChan <- false
				}
			}

			if checkNum == 0 {
				if cctx.Err() == nil {
					passChan <- true
				}
				return
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
					return
				}
				if err != nil && cctx.Err() == nil {
					errGot <- err
					return
				}
			case <-ctx.Done():
				if cctx.Err() == nil {
					errGot <- ctx.Err()
				}
				return
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
