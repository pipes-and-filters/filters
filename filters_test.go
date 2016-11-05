package filters

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type fakeReaderWriter struct {
}

func (f fakeReaderWriter) Read(p []byte) (n int, err error) {
	spew.Dump("read called")
	fmt.Println(string(p))
	return len(p), nil
}

func (f fakeReaderWriter) Write(p []byte) (n int, err error) {
	spew.Dump("write called")
	fmt.Println(string(p))
	return len(p), nil
}

func TestNew(t *testing.T) {
	f := New()
	if f.Arguments == nil {
		t.Error("Expected non nil slice")
	}
}

func TestRun(t *testing.T) {
	f := New()
	f.Command = "ls"
	frw := fakeReaderWriter{}
	e := f.Exec()
	e.SetOutput(frw)
	//bytes.NewReader([]byte(""))
	err := e.Run()
	if err != nil {
		t.Error(err)
	}

}

func TestChainRun(t *testing.T) {
	f := New()
	f.Command = "ls"
	fg := New()
	fg.Command = "grep"
	fg.Arguments = []string{"filter"}
	fx := New()
	fx.Command = "xargs"
	fx.Arguments = []string{"-n", "4"}
	var frw bytes.Buffer
	c := Chain{Filters: []Filter{f, fg, fx}}
	e, err := c.Exec()
	if err != nil {
		t.Error(err)
	}
	e.SetOutput(&frw)
	err = e.Run()
	if err != nil {
		t.Error(err)
	}
	io.Copy(os.Stdout, &frw)

}

func TestChainFromFile(t *testing.T) {
	c, err := ChainFile("chain.yml")
	if err != nil {
		t.Error(err)
	}
	spew.Dump(c)
}

func TestExecuteChainFromFile(t *testing.T) {
	c, err := ChainFile("chain.yml")
	if err != nil {
		t.Error(err)
	}
	var frw bytes.Buffer
	e, err := c.Exec()
	if err != nil {
		t.Error(err)
	}
	e.SetOutput(&frw)
	err = e.Run()
	if err != nil {
		t.Error(err)
	}
	io.Copy(os.Stdout, &frw)
}

func TestChainsFromFile(t *testing.T) {
	c, err := ChainsFile("chains.yml")
	if err != nil {
		t.Error(err)
	}
	spew.Dump(c)
}
