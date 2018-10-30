package check

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestCheck(t *testing.T) {
	var someError = errors.New("some error")

	gError := func(pass bool) error {
		ms := rand.Int63n(10) * 10
		time.Sleep(time.Millisecond * time.Duration(ms))
		if pass {
			return nil
		}
		return someError
	}

	gBool := func(pass bool) bool {
		return gError(pass) == nil
	}

	gCheckpoints := func(pass bool) []CheckFunc {
		checkPoints := []CheckFunc{}
		for i := 0; i < 10; i++ {
			var fn CheckFunc = func() bool {
				return gBool(pass)
			}
			checkPoints = append(checkPoints, fn)
		}
		return checkPoints
	}

	gErrCheckpoints := func(pass bool) []CheckErrFunc {
		checkPoints := []CheckErrFunc{}
		for i := 0; i < 10; i++ {
			var fn CheckErrFunc = func() error {
				return gError(pass)
			}
			checkPoints = append(checkPoints, fn)
		}
		return checkPoints
	}

	t.Run("pass", func(t *testing.T) {
		checkPoints := gCheckpoints(true)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if !Passed(ctx, checkPoints) {
			t.Error("not going to happen")
		}
	})

	t.Run("not pass", func(t *testing.T) {
		checkPoints := gCheckpoints(false)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if Passed(ctx, checkPoints) {
			t.Error("not going to happen")
		}
	})

	t.Run("error", func(t *testing.T) {
		checkPoints := gErrCheckpoints(true)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if NoError(ctx, checkPoints) != nil {
			t.Error("not going to happen")
		}
	})

	t.Run("no error", func(t *testing.T) {
		checkPoints := gErrCheckpoints(false)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if NoError(ctx, checkPoints) == nil {
			t.Error("not going to happen")
		}
	})

	t.Run("not pass with all cancel", func(t *testing.T) {
		checkPoints := gCheckpoints(false)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		if Passed(ctx, checkPoints) {
			t.Error("not going to happen")
		}
	})

	t.Run("error with all cancel", func(t *testing.T) {
		checkPoints := gErrCheckpoints(true)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		if NoError(ctx, checkPoints) != nil {
			t.Error("not going to happen")
		}
	})

	t.Run("not pass with some cancel", func(t *testing.T) {
		checkPoints := gCheckpoints(false)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		defer cancel()

		if Passed(ctx, checkPoints) {
			t.Error("not going to happen")
		}
	})

	t.Run("error with some cancel", func(t *testing.T) {
		checkPoints := gErrCheckpoints(true)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		defer cancel()

		if NoError(ctx, checkPoints) != nil {
			t.Error("not going to happen")
		}
	})
}
