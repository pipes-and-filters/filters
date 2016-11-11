package filters

import (
	"bytes"
	"fmt"

	"log"
)

func Example_simple() {
	text := "hello world"
	filter := New()
	filter.Command = "echo"
	filter.Argument(text)
	var buf bytes.Buffer
	exec := filter.Exec()
	exec.SetOutput(&buf)
	err := exec.Run()
	if err != nil {
		log.Print(err)
	}
	fmt.Println(buf.String())
	// Output:
	// hello world
}

func Example_input() {
	text := []byte("hello world\ngoodbye world")
	filter := New()
	filter.Command = "grep"
	filter.Argument("goodbye")
	var buf bytes.Buffer
	exec := filter.Exec()
	exec.SetInput(bytes.NewReader(text))
	exec.SetOutput(&buf)
	err := exec.Run()
	if err != nil {
		log.Print(err)
	}
	fmt.Println(buf.String())
	// Output:
	// goodbye world
}
