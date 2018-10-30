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
	printf("sleep %s", d)
	time.Sleep(d)
	if time.Now().Nanosecond()%4 == 0 && *noError == false {
		printf("error sleep %s", d)
		return fmt.Errorf("error: sleep %s", d)
	}
	printf("done sleep %s", d)
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
	checkErrA := func() bool { return generateErr(time.Second) == nil }
	checkErrB := func() bool { return generateErr(time.Second/2) == nil }
	checkErrC := func() bool { return generateErr(time.Second/3) == nil }
	checkErrD := func() bool {
		return generateErr(time.Second/4-100*time.Microsecond) == nil
	}
	checkErrE := func() bool { return generateErr(time.Second/4) == nil }

	checkPoints := []check.CheckFunc{
		checkErrA,
		checkErrB,
		checkErrC,
		checkErrD,
		checkErrE,
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
	checkErrA := func() error { return generateErr(time.Second) }
	checkErrB := func() error { return generateErr(time.Second / 2) }
	checkErrC := func() error { return generateErr(time.Second / 3) }
	checkErrD := func() error {
		return generateErr(time.Second/4 - 100*time.Microsecond)
	}
	checkErrE := func() error { return generateErr(time.Second / 4) }

	checkPoints := []check.CheckErrFunc{
		checkErrA,
		checkErrB,
		checkErrC,
		checkErrD,
		checkErrE,
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
	time.Sleep(time.Second)

	println("==testCheckError==")
	testCheckError()

	println("==end==")
	time.Sleep(time.Second * 2)
}
