[![codecov](https://codecov.io/gh/nikolaydubina/read-seek-peeker/graph/badge.svg?token=dWs1oSWSRU)](https://codecov.io/gh/nikolaydubina/read-seek-peeker)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/read-seek-peeker)](https://goreportcard.com/report/github.com/nikolaydubina/read-seek-peeker)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/read-seek-peeker#section-readme.svg)](https://pkg.go.dev/github.com/nikolaydubina/read-seek-peeker#section-readme)

```go
type ReadSeekPeeker interface {
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Peek(n int) ([]byte, error)
}
```

```go
s := "hello beautiful wonderful world!"
r := readseekpeeker.NewBufferedReadSeekPeeker(strings.NewReader(s), 5)

b := make([]byte, 5)
r.Read(b)

peek, _ := r.Peek(11)      // 😳

r.Seek(21, io.SeekCurrent) // 😮

rest, _ := io.ReadAll(r)

fmt.Println(string(b), string(peek), string(rest))
// Output: hello  beautiful  world!
```

As of `2024-02-22`, standard go packages allow either:  
A) `Seek()` by [io.ReadSeeker](https://pkg.go.dev/io#ReadSeeker)  
B) `Peek()` by [bufio.Reader](https://pkg.go.dev/bufio#Reader.Peek)  
..but not both! This package adds that.
