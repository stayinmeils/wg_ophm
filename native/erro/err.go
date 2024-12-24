package erro

import (
	"errors"
	"fmt"
	"os"
)

var Err chan error
var File *os.File
var TestFunc func() error
var Count int
var Buf []byte

var Count4 int
var Count6 int

func Errinit(fd int) {
	Err = make(chan error)
	Count = 0
	Count6 = 0
	Count4 = 0
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
