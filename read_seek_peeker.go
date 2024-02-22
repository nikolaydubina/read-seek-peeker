package readseekpeeker

import "io"

type ReadSeekPeeker interface {
	io.ReadSeeker
	Peek(n int) ([]byte, error)
}
