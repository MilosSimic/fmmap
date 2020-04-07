package fmmap

import (
	"errors"
	"os"
	"sync"
	"syscall"
)

type FMMAP struct {
	data []byte
	fd   int
	file *os.File
	mu   sync.Mutex
}

func (fmmap *FMMAP) Close() error {
	fmmap.file.Close()
	err := syscall.Munmap(fmmap.data)
	if err != nil {
		return err
	}
	return nil
}

func (fmmap *FMMAP) Update(data []byte) error {
	fmmap.mu.Lock()
	defer fmmap.mu.Unlock()

	if len(data) != len(fmmap.data) {
		err := fmmap.ftruncate(len(data))
		if err != nil {
			return err
		}
		err = fmmap.mmap()
		if err != nil {
			return err
		}
	}
	copy(fmmap.data, data)
	return nil
}

func (fmmap *FMMAP) UpdateFrom(i int, data []byte) error {
	fmmap.mu.Lock()
	defer fmmap.mu.Unlock()

	if len(fmmap.data) < i {
		return errors.New("Current file do not containt that starting position")
	}
	copy(fmmap.data[i:], data)
	err := fmmap.mmap()
	if err != nil {
		return err
	}
	return nil
}

func (fmmap *FMMAP) UpdateTo(i int, data []byte) error {
	fmmap.mu.Lock()
	defer fmmap.mu.Unlock()

	if len(fmmap.data) < i {
		return errors.New("Current file do not containt that starting position")
	}
	copy(fmmap.data[:i], data)
	err := fmmap.mmap()
	if err != nil {
		return err
	}
	return nil
}

func (fmmap *FMMAP) UpdateRange(i, j int, data []byte) error {
	fmmap.mu.Lock()
	defer fmmap.mu.Unlock()

	if len(fmmap.data) < i || len(fmmap.data) < j {
		return errors.New("Current file do not containt that starting position")
	}
	copy(fmmap.data[i:j], data)
	err := fmmap.mmap()
	if err != nil {
		return err
	}
	return nil
}

func (fmmap *FMMAP) Get() []byte {
	fmmap.mu.Lock()
	defer fmmap.mu.Unlock()

	return fmmap.data
}

func (fmmap *FMMAP) GetFrom(i int) []byte {
	fmmap.mu.Lock()
	defer fmmap.mu.Unlock()

	return fmmap.data[i:]
}

func (fmmap *FMMAP) GetTo(i int) []byte {
	fmmap.mu.Lock()
	defer fmmap.mu.Unlock()

	return fmmap.data[:i]
}

func (fmmap *FMMAP) GetRange(i, j int) []byte {
	fmmap.mu.Lock()
	defer fmmap.mu.Unlock()

	return fmmap.data[i:j]
}

func (fmmap *FMMAP) GetFile() *os.File {
	return fmmap.file
}

func (fmmap *FMMAP) open(filename string, flags int) {
	f, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	fmmap.file = f
	fmmap.fd = int(f.Fd())
}

func (fmmap *FMMAP) mmap() error {
	f, err := fmmap.file.Stat()
	if err != nil {
		return err
	}

	if int(f.Size()) != 0 {
		data, err := syscall.Mmap(fmmap.fd, 0, int(f.Size()), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			return err
		}
		fmmap.data = data
	}
	return nil
}

func (fmmap *FMMAP) ftruncate(size int) error {
	err := syscall.Ftruncate(fmmap.fd, int64(size))
	if err != nil {
		return err
	}
	return nil
}

func NewFile(file string, flags int) (*FMMAP, error) {
	fmmap := &FMMAP{}
	fmmap.open(file, flags)
	err := fmmap.mmap()
	if err != nil {
		return nil, err
	}
	return fmmap, nil
}
