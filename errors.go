package filters

import "fmt"

// Error is a filter error
type Error struct {
	command string
	err     []byte
}

// Write implements Writer interface
func (e *Error) Write(p []byte) (n int, err error) {
	e.err = append(e.err, p...)
	return len(p), nil
}

// Error implements Error interface
func (e Error) Error() string {
	return fmt.Sprintf("%v: %q", e.command, e.err)
}
