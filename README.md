# taskflow

[![GoDoc](https://godoc.org/github.com/cfsghost/taskflow?status.svg)](http://godoc.org/github.com/cfsghost/taskflow)

Golang library to build task flow architecture.

## Benchmark

Here is result of benchmark with testing tools:

```shell
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/cfsghost/taskflow
BenchmarkSingleTask
BenchmarkSingleTask-16    	 5083108	       213 ns/op
BenchmarkTwoTasks
BenchmarkTwoTasks-16      	 2585607	       458 ns/op
PASS
ok  	github.com/cfsghost/taskflow	3.195s
```

## License
Licensed under the MIT License

## Authors
Copyright(c) 2021 Fred Chien <cfsghost@gmail.com>
