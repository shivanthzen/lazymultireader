package lazymultireader

import (
	"io"
	"strings"
	"testing"
)

type stringReadOpener struct {
	strings.Reader
}

func (b *stringReadOpener) Open() error {
	// nop

	return nil
}

func NewStringReadOpener(s string) *stringReadOpener {
	return &stringReadOpener{
		*strings.NewReader(s),
	}
}

func Test_eofReader_Read(t *testing.T) {
	type args struct {
		in0 []byte
	}
	tests := []struct {
		name    string
		e       eofReader
		args    args
		want    int
		wantErr bool
	}{
		{
			"Should return EOF",
			eofReader{},
			args{nil},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := eofReader{}
			got, err := e.Read(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("eofReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("eofReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiReader(t *testing.T) {
	var mr io.Reader
	var buf []byte
	nread := 0
	withFooBar := func(tests func()) {
		r1 := NewStringReadOpener("foo ")
		r2 := NewStringReadOpener("")
		r3 := NewStringReadOpener("bar")
		mr = NewLazyMultiReader(r1, r2, r3)
		buf = make([]byte, 20)
		tests()
	}
	expectRead := func(size int, expected string, eerr error) {
		nread++
		n, gerr := mr.Read(buf[0:size])
		if n != len(expected) {
			t.Errorf("#%d, expected %d bytes; got %d",
				nread, len(expected), n)
		}
		got := string(buf[0:n])
		if got != expected {
			t.Errorf("#%d, expected %q; got %q",
				nread, expected, got)
		}
		if gerr != eerr {
			t.Errorf("#%d, expected error %v; got %v",
				nread, eerr, gerr)
		}
		buf = buf[n:]
	}
	withFooBar(func() {
		expectRead(2, "fo", nil)
		expectRead(5, "o ", nil)
		expectRead(5, "bar", nil)
		expectRead(5, "", io.EOF)
	})
	withFooBar(func() {
		expectRead(4, "foo ", nil)
		expectRead(1, "b", nil)
		expectRead(3, "ar", nil)
		expectRead(1, "", io.EOF)
	})
	withFooBar(func() {
		expectRead(5, "foo ", nil)
	})
}
