package filters

import (
	"fmt"
	"io"
	"testing"
)

type fakeReaderWriter struct {
	Input  []byte
	Output []byte
	i      int64
}

func (f *fakeReaderWriter) Read(p []byte) (n int, err error) {
	if f.i >= int64(len(f.Input)) {
		return 0, io.EOF
	}
	n = copy(p, f.Input[f.i:])
	f.i += int64(n)
	return
}

func (f *fakeReaderWriter) Write(p []byte) (n int, err error) {
	f.Output = append(f.Output, p...)
	return len(p), nil
}

func TestNew(t *testing.T) {
	f := New()
	if f.Arguments == nil {
		t.Error("Expected non nil slice")
	}
}

func TestFilterFile(t *testing.T) {
	_, err := FilterFile("filter.yml")
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	ts := "hello world"
	f := New()
	f.Command = "echo"
	f.Argument(ts)
	frw := fakeReaderWriter{}
	e := f.Exec()
	e.SetOutput(&frw)
	err := e.Run()
	if err != nil {
		t.Error(err)
	}
	if string(frw.Output) != fmt.Sprintf("%v\n", ts) {
		t.Errorf("Expected ts with carriage return received %v", string(frw.Output))
	}

}

func TestChainRun(t *testing.T) {
	ts := "Alphabet city is haunted\nConstantina feels right at home\n"
	ts += "She probably wont say youre wrong\n"
	ts += "Youre already wrong\nYoure already wrong\n"
	f := New()
	f.Command = "cat"
	fg := New()
	fg.Command = "grep"
	fg.Argument("wrong")
	fx := New()
	fx.Command = "xargs"
	fx.Argument("-n")
	fx.Argument("3")
	frw := fakeReaderWriter{Input: []byte(ts)}
	//c := Chain{Filters: []Filter{f, fg, fx}}
	c := Chain{Filters: []Filter{f, fg, fx}}
	e, err := c.Exec()
	if err != nil {
		t.Error(err)
	}
	e.SetInput(&frw)
	e.SetOutput(&frw)
	err = e.Run()
	if err != nil {
		t.Error(err)
	}
	exp := "She probably wont\nsay youre wrong\n"
	exp += "Youre already wrong\nYoure already wrong\n"
	if string(frw.Output) != fmt.Sprintf("%v", exp) {
		t.Errorf("Expected: \n%v \ngot \n%v", exp, string(frw.Output))
	}

}

func TestChainFromFile(t *testing.T) {
	_, err := ChainFile("chain.yml")
	if err != nil {
		t.Error(err)
	}
}

func TestExecuteChainFromFile(t *testing.T) {
	ts := "Alphabet city is haunted\nConstantina feels right at home\n"
	ts += "She probably wont say youre wrong\n"
	ts += "Youre already wrong\nYoure already wrong\n"
	c, err := ChainFile("chain.yml")
	if err != nil {
		t.Error(err)
	}
	frw := fakeReaderWriter{Input: []byte(ts)}
	e, err := c.Exec()
	if err != nil {
		t.Error(err)
	}
	e.SetInput(&frw)
	e.SetOutput(&frw)
	err = e.Run()
	if err != nil {
		t.Error(err)
	}
	exp := "She probably wont\nsay youre wrong\n"
	exp += "Youre already wrong\nYoure already wrong\n"
	if string(frw.Output) != fmt.Sprintf("%v", exp) {
		t.Errorf("Expected: \n%v \ngot \n%v", exp, string(frw.Output))
	}
}

func TestChainsFromFile(t *testing.T) {
	_, err := ChainsFile("chains.yml")
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	c, err := ChainsFile("chains.yml")
	if err != nil {
		t.Error(err)
	}
	_, err = c.Get("FirstChain")
	if err != nil {
		t.Error(err)
	}
	_, err = c.Get("boguschain")
	if err != ErrChainDoesNotExist {
		t.Error("Expected chain does not exist error")
	}
}

func TestFromFile(t *testing.T) {
	file := "bogus.yml"
	c := NewChains()
	err := fromFile(file, c)
	if err != ErrNotPointer {
		t.Error("Expected not pointer error")
	}
	err = fromFile(file, &c)
	if err.Error() != "stat bogus.yml: no such file or directory" {
		t.Error("Expected bad file error")
	}

}
