package main

import (
	"errors"
	"os"
)

var testData = []byte("0")

func openFile(file string, flags int) (*os.File, error) {
	if len(file) <= 0 {
		return nil, errors.New("file must have a name")
	}

	f, err := os.OpenFile(file, flags, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func stat(f *os.File) (bool, error) {
	fi, err := f.Stat()
	if err != nil {
		return false, err
	}

	if fi.Size() <= 0 {
		return true, nil
	}
	return false, nil
}
