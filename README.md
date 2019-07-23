# Check

[![GoDoc](https://godoc.org/github.com/wrfly/check?status.svg)](https://godoc.org/github.com/wrfly/check)

Check multiple errors/functions

If you want to check some functions currently, and lazy to write some code,
you can use this lib.

This lib will check the checkpoints you provided, run then in currency, and
return the first error or false if found. Otherwise, return true or nil error,
means all checkpoints are passed.

Get this package and run the example:

```bash
go get github.com/wrfly/check

# generate errors or false randomly
go run $GOPATH/src/github.com/wrfly/check/example/main.go

# no error or false
go run $GOPATH/src/github.com/wrfly/check/example/main.go -ne
```

## Go test result

```txt
go test -v --race --cover .
=== RUN   TestCheck
=== RUN   TestCheck/pass
=== RUN   TestCheck/not_pass
=== RUN   TestCheck/error
=== RUN   TestCheck/no_error
=== RUN   TestCheck/not_pass_with_all_cancel
=== RUN   TestCheck/error_with_all_cancel
=== RUN   TestCheck/not_pass_with_some_cancel
=== RUN   TestCheck/error_with_some_cancel
=== RUN   TestCheck/no_check_points
--- PASS: TestCheck (0.27s)
    --- PASS: TestCheck/pass (0.09s)
    --- PASS: TestCheck/not_pass (0.01s)
    --- PASS: TestCheck/error (0.10s)
    --- PASS: TestCheck/no_error (0.00s)
    --- PASS: TestCheck/not_pass_with_all_cancel (0.00s)
    --- PASS: TestCheck/error_with_all_cancel (0.00s)
    --- PASS: TestCheck/not_pass_with_some_cancel (0.00s)
    --- PASS: TestCheck/error_with_some_cancel (0.05s)
    --- PASS: TestCheck/no_check_points (0.00s)
PASS
coverage: 100.0% of statements
ok      github.com/wrfly/check  1.281s  coverage: 100.0% of statements
```
