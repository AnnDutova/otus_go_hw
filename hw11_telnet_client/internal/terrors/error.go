package terrors

import "fmt"

type ClientErr struct {
	Msg string
	Err error
}

func (e ClientErr) Error() string {
	return fmt.Sprintf("%s: %s\n", e.Msg, e.Err)
}
