package fmmap

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

type FMMAP struct {
	data []byte
	fd   int
	file *os.File
}

func (fmmap *FMMAP) Close() {
	fmmap.file.Close()
	err := syscall.Munmap(fmmap.data)
	if err != nil {
		fmt.Println("Error unmap data: ", err)
	}
}

func (fmmap *FMMAP) Update(data []byte) error {
	if len(data) != len(fmmap.data) {
		err := fmmap.ftruncate(len(data))
		if err != nil {
			return err
		}
	}
	copy(fmmap.data, data)
	fmmap.mmap()
	err := fmmap.msync()
	if err != nil {
		return nil
	}

	return nil
}

func (fmmap *FMMAP) UpdateFrom(i int, data []byte) error {
	if len(fmmap.data) < i {
		return errors.New("Current file do not containt that starting position")
	}
	copy(fmmap.data[i:], data)
	fmmap.mmap()

	return nil
}

func (fmmap *FMMAP) UpdateTo(i int, data []byte) error {
	if len(fmmap.data) < i {
		return errors.New("Current file do not containt that starting position")
	}
	copy(fmmap.data[:i], data)
	fmmap.mmap()

	return nil
}

func (fmmap *FMMAP) Updaterange(i, j int, data []byte) error {
	if len(fmmap.data) < i || len(fmmap.data) < j {
		return errors.New("Current file do not containt that starting position")
	}
	copy(fmmap.data[i:j], data)
	fmmap.mmap()

	return nil
}

func (fmmap *FMMAP) Get() []byte {
	return fmmap.data
}

func (fmmap *FMMAP) GetFrom(i int) []byte {
	return fmmap.data[i:]
}

func (fmmap *FMMAP) GetTo(i int) []byte {
	return fmmap.data[:i]
}

func (fmmap *FMMAP) GetRange(i, j int) []byte {
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
	// fmmap.mmap()
}

func (fmmap *FMMAP) mmap() {
	f, err := fmmap.file.Stat()
	if err != nil {
		fmt.Println("Could not stat file: ", err)
		return
	}

	data, err := syscall.Mmap(fmmap.fd, 0, int(f.Size()), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("Error mmapping: ", err)
		return
	}
	fmmap.data = data
}

func (fmmap *FMMAP) msync() error {
	return unix.Msync(fmmap.data, unix.MS_SYNC)
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
	fmmap.mmap()

	return fmmap, nil
}
