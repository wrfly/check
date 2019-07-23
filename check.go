package check

import (
	"context"
)

// Func is the check function that only returns true or false
type Func func() bool

// FuncWithErr returns error if any
type FuncWithErr func() error

type e struct{}

func wrapCheck(ctx context.Context, start chan e, cf Func, noErrChan chan e) chan e {
	failed := make(chan e)
	go func() {
		c := make(chan bool, 1)
		defer close(c)
		<-start

		select {
		case <-ctx.Done():
			close(failed)
		case c <- cf():
			if !<-c {
				close(failed)
			} else {
				noErrChan <- e{}
			}
		}
	}()
	return failed
}

func wrapCheckWithError(ctx context.Context, start chan e, cf FuncWithErr) chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		<-start

		select {
		case <-ctx.Done():
		case ch <- cf():
		}
	}()
	return ch
}

// Passed returns true if all check points passed, otherwise, returns false
func Passed(ctx context.Context, checkPoints []Func) bool {
	checkNum := len(checkPoints)
	if checkNum == 0 {
		return true
	}
	passedChan := make(chan e, checkNum)
	failedChan := make(chan e, checkNum)

	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	start := make(chan e)
	for _, cf := range checkPoints {
		go mergeE(failedChan, wrapCheck(cCtx, start, cf, passedChan))
	}

	passChan := make(chan bool)

	go func() {
		defer close(passChan)
		defer close(failedChan)
		defer close(passedChan)
		for {
			if checkNum == 0 {
				// all checkpoints passed, if the check context
				// is still on going (within the context)
				// then the check passed, return to close all the channels
				if cCtx.Err() == nil {
					passChan <- true
				}
				return
			}

			select {
			case <-passedChan:
				// one check passed
				checkNum--

			case <-failedChan:
				checkNum--
				// one check failed, then if the
				// check context is still on going,
				// then the check list is failed
				if cCtx.Err() == nil {
					passChan <- false
				}

			case <-ctx.Done():
				passChan <- false
			}

		}
	}()

	close(start) // start all checks at the same time

	return <-passChan
}

// NoError returns the first error it got, if all passed, returns nil
func NoError(ctx context.Context, checkPoints []FuncWithErr) error {
	checkNum := len(checkPoints)
	if checkNum == 0 {
		return nil
	}
	errorChan := make(chan error, len(checkPoints))

	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	start := make(chan e)
	for _, cf := range checkPoints {
		go mergeErr(errorChan, wrapCheckWithError(cCtx, start, cf))
	}

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
				if err != nil && cCtx.Err() == nil {
					errGot <- err
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	close(start)

	return <-errGot
}

func mergeE(dest, src chan e) {
	dest <- <-src
}

func mergeErr(dest, src chan error) {
	dest <- <-src
}
