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