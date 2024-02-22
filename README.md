[![codecov](https://codecov.io/gh/nikolaydubina/readseekpeeker/graph/badge.svg?token=dWs1oSWSRU)](https://codecov.io/gh/nikolaydubina/readseekpeeker)
https://pkg.go.dev/github.com/nikolaydubina/readseekpeeker
[![Go Reference](https://pkg.go.dev/github.com/nikolaydubina/readseekpeeker.svg)](https://pkg.go.dev/github.com/nikolaydubina/readseekpeeker)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/fpdecimal)](https://goreportcard.com/report/github.com/nikolaydubina/fpdecimal)

```go
type ReadSeekPeeker interface {
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Peek(n int) ([]byte, error)
}
```

## Why?

As of `2024-02-22`, standard go packages allow either:  
A) `Seek()` by [io.ReadSeeker](https://pkg.go.dev/io#ReadSeeker)  
B) `Peek()` by [bufio.Reader](https://pkg.go.dev/bufio#Reader.Peek)  

but not both! This package adds that.
