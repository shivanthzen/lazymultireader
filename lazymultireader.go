package lazymultireader

import (
	"io"
)

type eofReader struct{}

func (eofReader) Read([]byte) (int, error) {
	return 0, io.EOF
}

type ReadOpener interface {
	Open() error
	io.Reader
}

type LazyMultiReader struct {
	reader  io.Reader
	readers []ReadOpener
}

func NewLazyMultiReader(s ...ReadOpener) io.Reader {
	if s[0].Open() != nil {
		return nil
	}
	return &LazyMultiReader{
		s[0],
		s,
	}
}

func (mr *LazyMultiReader) Read(p []byte) (n int, err error) {
	if mr.reader == nil {
		return 0, io.EOF
	}
	for len(mr.readers) > 0 {
		n, err = mr.reader.Read(p)
		if err == io.EOF {
			mr.reader = eofReader{} // permit earlier GC
			mr.readers = mr.readers[1:]
			if len(mr.readers) < 1 {
				return
			}
			reader := mr.readers[0]
			err := reader.Open()
			if err != nil {
				return 0, err
			}
			mr.reader = reader
		}
		if n > 0 || err != io.EOF {
			if err == io.EOF && len(mr.readers) > 0 {
				err = nil
			}
			return
		}
	}
	return 0, io.EOF
}
