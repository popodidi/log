package file

import (
	"io"
	"os"
	"path"
	"sync"
)

func newSingle(filename string) *single {
	return &single{
		filename: filename,
	}
}

var _ io.WriteCloser = (*single)(nil)

type single struct {
	sync.Once // to open file
	filename  string
	f         *os.File
}

func (s *single) Write(p []byte) (n int, err error) {
	s.Do(func() { err = s.open() })
	if err != nil {
		return
	}
	return s.f.Write(p)
}

func (s *single) open() error {
	// make sure the dir exists
	dir := path.Dir(s.filename)
	dirExist, err := exist(dir)
	if err != nil {
		return err
	}
	if !dirExist {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	// maker sure the file exists
	fileExist, err := exist(s.filename)
	if err != nil {
		return err
	}
	if fileExist {
		s.f, err = os.OpenFile( // nolint: gosec
			s.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	} else {
		s.f, err = os.Create(s.filename)
	}
	return err
}

func (s *single) exist() (bool, error) {
	return exist(s.filename)
}

func (s *single) Close() error {
	if s.f == nil {
		return nil
	}
	return s.f.Close()
}
