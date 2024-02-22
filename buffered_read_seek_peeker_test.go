package readseekpeeker_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	readseekpeeker "github.com/nikolaydubina/read-seek-peeker"
)

func ExampleBufferedReadSeekPeeker() {
	s := "hello beautiful wonderful world!"
	r := readseekpeeker.NewBufferedReadSeekPeeker(strings.NewReader(s), 5)

	b := make([]byte, 5)
	r.Read(b)

	peek, _ := r.Peek(11)

	r.Seek(21, io.SeekCurrent)

	rest, _ := io.ReadAll(r)

	fmt.Println(string(b), string(peek), string(rest))
	// Output: hello  beautiful  world!
}

func TestBufferedReadSeekPeeker(t *testing.T) {
	r := readseekpeeker.NewBufferedReadSeekPeeker(bytes.NewReader([]byte("hi world !!! how you doing ?")), 5)

	func() {
		b := make([]byte, 3)
		n, err := r.Read(b)
		if err != nil {
			t.Fatal(err)
		}
		if n != 3 {
			t.Fatalf("expected 3, got %d", n)
		}
		if string(b) != "hi " {
			t.Fatalf("expected 'hi ', got '%s'", b)
		}
	}()

	func() {
		b, err := r.Peek(5)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "world" {
			t.Fatalf("expected 'world', got '%s'", b)
		}
	}()

	func() {
		n, err := r.Seek(3, io.SeekCurrent)
		if err != nil {
			t.Fatal(err)
		}
		if n != 6 {
			t.Fatalf("expected 6, got %d", n)
		}
	}()

	func() {
		b := make([]byte, 3)
		n, err := r.Read(b)
		if err != nil {
			t.Fatal(err)
		}
		if n != 3 {
			t.Fatalf("expected 3, got %d", n)
		}
		if string(b) != "ld " {
			t.Fatalf("expected 'ld ', got '%s'", b)
		}
	}()

	func() {
		b, err := r.Peek(10)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "!!! how yo" {
			t.Fatalf("expected '!!! how yo', got '%s'", b)
		}
	}()

	func() {
		b, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "!!! how you doing ?" {
			t.Fatalf("expected '!!! how you doing ?', got '%s'", b)
		}
	}()
}

func TestBufferedReadSeekPeeker_ResetBufferOnLargeSeeks(t *testing.T) {
	s := strings.Repeat("abcdefg", 8*1<<20/7) // 8 MB

	r := readseekpeeker.NewBufferedReadSeekPeeker(strings.NewReader(s), 5)

	n, err := r.Seek(2*7*1<<17, io.SeekCurrent) // 130 KB
	if err != nil {
		t.Fatal(err)
	}
	if n != 2*7*1<<17 {
		t.Fatalf("expected %d, got %d", 2*7*1<<17, n)
	}

	func() {
		b := make([]byte, 7)
		n, err := r.Read(b)
		if err != nil {
			t.Fatal(err)
		}
		if n != 7 {
			t.Fatalf("expected 7, got %d", n)
		}
		if string(b) != "abcdefg" {
			t.Fatalf("expected 'abcdefg', got '%s'", b)
		}
	}()

	func() {
		b, err := r.Peek(14)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "abcdefgabcdefg" {
			t.Fatalf("expected 'abcdefgabcdefg', got '%s'", b)
		}
	}()

	func() {
		n, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			t.Fatal(err)
		}
		if n != 2*7*1<<17+7 {
			t.Fatalf("expected %d, got %d", 2*7*1<<17+7, n)
		}
	}()
}

func TestBufferedReadSeekPeeker_Unsupported(t *testing.T) {
	r := readseekpeeker.NewBufferedReadSeekPeeker(strings.NewReader("abcdefg"), 5)

	func() {
		n, err := r.Seek(-1, io.SeekCurrent)
		if err == nil {
			t.Fatal("err expected")
		}
		if n != 0 {
			t.Fatalf("wrong offset")
		}
	}()

	func() {
		n, err := r.Seek(3, io.SeekStart)
		if err == nil {
			t.Fatal("err expected")
		}
		if n != 0 {
			t.Fatalf("wrong offset")
		}
	}()

	func() {
		n, err := r.Seek(3, io.SeekEnd)
		if err == nil {
			t.Fatal("err expected")
		}
		if n != 0 {
			t.Fatalf("wrong offset")
		}
	}()
}

type failingSeeker struct{}

func (s *failingSeeker) Read(p []byte) (n int, err error) { return 0, io.EOF }

func (s *failingSeeker) Seek(offset int64, whence int) (int64, error) { return 0, errors.New("err") }

func TestBufferedReadSeekPeeker_FailingSeeker(t *testing.T) {
	s := &failingSeeker{}
	r := readseekpeeker.NewBufferedReadSeekPeeker(s, 5)

	n, err := r.Seek(3, io.SeekCurrent)
	if err == nil {
		t.Fatal("err expected")
	}
	if n != 0 {
		t.Fatalf("wrong offset")
	}
}
