package check

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestCheck(t *testing.T) {
	gError := func(pass bool) error {
		ms := rand.Int63n(100) * 10
		time.Sleep(time.Millisecond * time.Duration(ms))
		if pass {
			return nil
		}
		return errors.New("")
	}

	gBool := func(pass bool) bool {
		return gError(pass) == nil
	}

	t.Run("pass", func(t *testing.T) {
		checkPoints := []CheckFunc{}
		for i := 0; i < 10; i++ {
			var fn CheckFunc = func() bool {
				return gBool(true)
			}
			checkPoints = append(checkPoints, fn)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if !Passed(ctx, checkPoints) {
			t.Error("not going to happen")
		}
	})

	t.Run("not pass", func(t *testing.T) {
		checkPoints := []CheckFunc{}
		for i := 0; i < 10; i++ {
			var fn CheckFunc = func() bool {
				return gBool(false)
			}
			checkPoints = append(checkPoints, fn)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if Passed(ctx, checkPoints) {
			t.Error("not going to happen")
		}
	})

	t.Run("error", func(t *testing.T) {
		checkPoints := []CheckErrFunc{}
		for i := 0; i < 10; i++ {
			var fn CheckErrFunc = func() error {
				return gError(true)
			}
			checkPoints = append(checkPoints, fn)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if NoError(ctx, checkPoints) != nil {
			t.Error("not going to happen")
		}
	})

	t.Run("no error", func(t *testing.T) {
		checkPoints := []CheckErrFunc{}
		for i := 0; i < 10; i++ {
			var fn CheckErrFunc = func() error {
				return gError(false)
			}
			checkPoints = append(checkPoints, fn)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if NoError(ctx, checkPoints) == nil {
			t.Error("not going to happen")
		}
	})
}
