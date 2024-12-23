package erro

import (
	"errors"
	"fmt"
	"os"
)

var Err chan error
var File *os.File
var TestFunc func() error

func Errinit(fd int) {
	Err = make(chan error)
	File = os.NewFile(uintptr(fd), "/dev/tun")
	count := 0
	TestFunc = func() error {
		buf := make([]byte, 65535)
		_, err := File.Read(buf)
		if err != nil {
			return errors.New(fmt.Sprintf("test filed %d", count) + err.Error())
		}
		count++
		return nil
	}
}
