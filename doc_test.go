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

func Example_chain() {
	text := "Alphabet city is haunted\n"
	text += "Constantina feels right at home\n"
	text += "She probably wont say youre wrong\n"
	text += "Youre already wrong\nYoure already wrong\n"
	filtercat := New()
	filtercat.Command = "cat"
	filtergrep := New()
	filtergrep.Command = "grep"
	filtergrep.Argument("wrong")
	filterxargs := New()
	filterxargs.Command = "xargs"
	filterxargs.Argument("-n")
	filterxargs.Argument("3")
	c := Chain{
		Filters: []Filter{
			filtercat,
			filtergrep,
			filterxargs,
		},
	}
	exec, err := c.Exec()
	if err != nil {
		log.Print(err)
	}
	exec.SetInput(bytes.NewReader([]byte(text)))
	var buf bytes.Buffer
	exec.SetOutput(&buf)
	err = exec.Run()
	if err != nil {
		log.Print(err)
	}
	fmt.Print(buf.String())
	// Output:
	// She probably wont
	// say youre wrong
	// Youre already wrong
	// Youre already wrong
}

func Example_chainfile() {
	text := "Alphabet city is haunted\n"
	text += "Constantina feels right at home\n"
	text += "She probably wont say youre wrong\n"
	text += "Youre already wrong\nYoure already wrong\n"
	c, err := ChainFile("chain.yml")
	if err != nil {
		log.Print(err)
	}
	exec, err := c.Exec()
	if err != nil {
		log.Print(err)
	}
	exec.SetInput(bytes.NewReader([]byte(text)))
	var buf bytes.Buffer
	exec.SetOutput(&buf)
	err = exec.Run()
	if err != nil {
		log.Print(err)
	}
	fmt.Print(buf.String())
	// Output:
	// She probably wont
	// say youre wrong
	// Youre already wrong
	// Youre already wrong
}
