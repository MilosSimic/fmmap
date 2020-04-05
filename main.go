package main

import (
	"fmt"
	"os"
)

func main() {
	d, err := NewFile("mmap_test", os.O_RDWR|os.O_CREATE)
	if err != nil {
		fmt.Println(err)
	}
	defer d.Close()

	fmt.Println(d)

	err = d.Update([]byte("12357567sfgdgh"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(d)
}
