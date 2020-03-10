package file

import (
	"io"
	"os"
	"path"
)

func newSingle(filename string) (*single, error) {
	dir := path.Dir(filename)
	// make sure the dir exists
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	if err != nil {
		return nil, err
	}
	// maker sure the file exists
	var f *os.File
	_, err = os.Stat(filename)
	if err != nil {
		f, err = os.Create(filename)
	} else {
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}
	if err != nil {
		return nil, err
	}
	return &single{
		f: f,
	}, nil
}

var _ io.WriteCloser = (*single)(nil)

type single struct {
	f *os.File
}

func (s *single) Write(p []byte) (n int, err error) {
	return s.f.Write(p)
}

func (s *single) Close() error {
	return s.f.Close()
}
