package erro

var Err chan error

func Errinit() {
	Err = make(chan error)
}
