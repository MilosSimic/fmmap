package fmmap

import (
	"errors"
	mmaplib "github.com/edsrzf/mmap-go"
	"os"
)

type FMMAP struct {
	f    *os.File
	mmap mmaplib.MMap
}

func (fmmap *FMMAP) Close() {
	fmmap.f.Close()
	fmmap.mmap.Unmap()
}

func (fmmap *FMMAP) Update(data []byte) error {
	fmmap.mmap = data
	fmmap.mmap.Flush()
	return nil
}

func (fmmap *FMMAP) UpdateFrom(i int, data []byte) error {
	if len(fmmap.mmap) < i {
		return errors.New("Current file do not containt that starting position")
	}

	copy(fmmap.mmap[i:], data)
	fmmap.mmap.Flush()
	return nil
}

func (fmmap *FMMAP) UpdateTo(i int, data []byte) error {
	if len(fmmap.mmap) < i {
		return errors.New("Current file do not containt that starting position")
	}

	copy(fmmap.mmap[:i], data)
	fmmap.mmap.Flush()
	return nil
}

func (fmmap *FMMAP) Updaterange(i, j int, data []byte) error {
	if len(fmmap.mmap) < i || len(fmmap.mmap) < j {
		return errors.New("Current file do not containt that starting position")
	}

	copy(fmmap.mmap[i:j], data)
	fmmap.mmap.Flush()
	return nil
}

func (fmmap *FMMAP) Get() []byte {
	return fmmap.mmap
}

func (fmmap *FMMAP) GetFrom(i int) []byte {
	return fmmap.mmap[i:]
}

func (fmmap *FMMAP) GetTo(i int) []byte {
	return fmmap.mmap[:i]
}

func (fmmap *FMMAP) GetRange(i, j int) []byte {
	return fmmap.mmap[i:j]
}

func (fmmap *FMMAP) GetFile() *os.File {
	return fmmap.f
}

func NewFile(file string, flags int) (*FMMAP, error) {
	f, err := openFile(file, flags)
	if err != nil {
		return nil, err
	}

	add, err := stat(f)
	if err != nil {
		return nil, err
	}

	if add {
		f.Write(testData)
	}

	mmap, err := mmaplib.Map(f, mmaplib.RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &FMMAP{
		f:    f,
		mmap: mmap,
	}, nil
}
