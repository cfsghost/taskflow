# taskflow

[![GoDoc](https://godoc.org/github.com/cfsghost/taskflow?status.svg)](http://godoc.org/github.com/cfsghost/taskflow)

Golang library to build task flow architecture.

## Benchmark

Here is result of benchmark with testing tools:

```shell
$ go test -bench=. -v
=== RUN   TestCreateTaskFlow
--- PASS: TestCreateTaskFlow (0.00s)
=== RUN   TestCreateEmptyTask
--- PASS: TestCreateEmptyTask (0.00s)
=== RUN   TestRemoveTask
--- PASS: TestRemoveTask (0.00s)
=== RUN   TestCreateCustomizedTask
--- PASS: TestCreateCustomizedTask (0.00s)
=== RUN   TestMultipleSend
--- PASS: TestMultipleSend (0.00s)
=== RUN   TestFanOutData
--- PASS: TestFanOutData (0.00s)
=== RUN   TestUnlink
--- PASS: TestUnlink (0.00s)
=== RUN   TestPrivateData
    taskflow_test.go:206: private data
--- PASS: TestPrivateData (0.00s)
goos: darwin
goarch: amd64
pkg: github.com/cfsghost/taskflow
BenchmarkSingleTask
BenchmarkSingleTask-16    	 5990356	       185 ns/op
BenchmarkTwoTasks
BenchmarkTwoTasks-16      	 2381484	       511 ns/op
BenchmarkTenTasks
BenchmarkTenTasks-16      	  617601	      2368 ns/op
PASS
ok  	github.com/cfsghost/taskflow	5.000s
```

## License
Licensed under the MIT License

## Authors
Copyright(c) 2021 Fred Chien <cfsghost@gmail.com>
