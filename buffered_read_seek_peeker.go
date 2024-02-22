package readseekpeeker

import (
	"bufio"
	"errors"
	"io"
)

// BufferedReadSeekPeeker can read, peek and seek by using buffered reader.
// Useful when you want to seek and peek at the same time.
type BufferedReadSeekPeeker struct {
	r      io.ReadSeeker
	b      *bufio.Reader
	offset int64
}

func NewBufferedReadSeekPeeker(r io.ReadSeeker, bufferSize int) *BufferedReadSeekPeeker {
	return &BufferedReadSeekPeeker{
		r:      r,
		b:      bufio.NewReaderSize(r, bufferSize),
		offset: 0,
	}
}

func (s *BufferedReadSeekPeeker) Peek(n int) ([]byte, error) { return s.b.Peek(n) }

func (s *BufferedReadSeekPeeker) Read(b []byte) (int, error) {
	// buffered reader may not advance underlying reader yet, recalculating offset
	n, err := s.b.Read(b)
	s.offset += int64(n)
	return n, err
}

func (s *BufferedReadSeekPeeker) Seek(n int64, whence int) (offset int64, err error) {
	if whence != io.SeekCurrent {
		return s.offset, errors.New("whence not supported, only io.SeekCurrent is supported")
	}
	if n < 0 {
		return s.offset, errors.New("negative offset not supported")
	}
	if n == 0 {
		return s.offset, nil
	}

	// if inside buffer, then discard prefix of buffer
	if n > 0 && n <= int64(s.b.Buffered()) {
		// as of 2024-02-22, Discard() is guaranteed to succeed if 0 <= n <= s.Buffered()
		s.b.Discard(int(n))
		s.offset += int64(n)
		return s.offset, nil
	}

	// move underlying seeker. buffered reader may advanced already, so sync them.
	if s.offset, err = s.r.Seek(n-int64(s.b.Buffered()), whence); err != nil {
		return s.offset, err
	}
	s.b.Reset(s.r)

	return s.offset, nil
}
