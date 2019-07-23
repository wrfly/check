package main

// run with `go run . -ne`
// or just `go run .`

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/wrfly/check"
)

var noError = flag.Bool("ne", false, "no error")

func generateErr(d time.Duration) error {
	returnError := time.Now().Nanosecond()%4 != 0 && *noError == false
	if returnError {
		printf("start to sleep %s with error", d)
	} else {
		printf("start to sleep %s", d)
	}
	time.Sleep(d)
	printf("done sleep %s", d)
	if returnError {
		return fmt.Errorf("sleep %s", d)
	}
	return nil
}

func printf(format string, msg interface{}) {
	fmt.Printf("[%s] %s\n",
		time.Now().Format("15:04:05.000"),
		fmt.Sprintf(format, msg))
}

func println(msg string) {
	printf("%s", msg)
}

func testCheck() {
	checkPoints := []check.Func{
		func() bool { return generateErr(time.Second) == nil },
		func() bool { return generateErr(time.Second/2) == nil },
		func() bool { return generateErr(time.Second/3) == nil },
		func() bool { return generateErr(time.Second/4-100*time.Microsecond) == nil },
		func() bool { return generateErr(time.Second/4) == nil },
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	passed := check.Passed(ctx, checkPoints)

	if passed {
		println("->> passed")
	} else {
		println("->> not passed")
	}

}

func testCheckError() {

	checkPoints := []check.FuncWithErr{
		func() error { return generateErr(time.Second) },
		func() error { return generateErr(time.Second / 2) },
		func() error { return generateErr(time.Second / 3) },
		func() error { return generateErr(time.Second/4 - 100*time.Microsecond) },
		func() error { return generateErr(time.Second / 4) },
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := check.NoError(ctx, checkPoints); err == nil {
		println("->> passed")
	} else {
		printf("->> got error: %s", err)
	}

}

func main() {
	flag.Parse()

	println("==testCheck==")
	testCheck()
	time.Sleep(time.Second * 1)

	println("==testCheckError==")
	testCheckError()

	println("==end==")

	time.Sleep(time.Second * 1)
}
