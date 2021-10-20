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
--- PASS: TestPrivateData (0.00s)
goos: darwin
goarch: amd64
pkg: github.com/cfsghost/taskflow
BenchmarkSingleTask
BenchmarkSingleTask-16                	 6052552	       193 ns/op
BenchmarkTwoTasks
BenchmarkTwoTasks-16                  	 2380485	       525 ns/op
BenchmarkTenTasks_4_Workers
BenchmarkTenTasks_4_Workers-16        	  646082	      2322 ns/op
BenchmarkTenTasks_8_Workers
BenchmarkTenTasks_8_Workers-16        	  679878	      1885 ns/op
BenchmarkHundredTasks_4_Workers
BenchmarkHundredTasks_4_Workers-16    	   36993	     35597 ns/op
BenchmarkHundredTasks_8_Workers
BenchmarkHundredTasks_8_Workers-16    	   48466	     27605 ns/op
PASS
ok  	github.com/cfsghost/taskflow	11.858s
```

## License
Licensed under the MIT License

## Authors
Copyright(c) 2021 Fred Chien <cfsghost@gmail.com>
