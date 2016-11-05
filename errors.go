package filters

import "fmt"

type Error struct {
	command string
	err     []byte
}

func (e *Error) Write(p []byte) (n int, err error) {
	e.err = append(e.err, p...)
	return len(p), nil
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %q", e.command, e.err)
}
