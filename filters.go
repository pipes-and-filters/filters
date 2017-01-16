package filters

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"syscall"

	yaml "gopkg.in/yaml.v2"
)

var (
	// ErrNotPointer the passed in vvariable was not a pointer
	ErrNotPointer = errors.New("Unmarshal requires pointer")
	// ErrChainDoesNotExist the chain was not found in the list of chains
	ErrChainDoesNotExist = errors.New("Chain does not exist")
)

// Filter is one command to be used as a filter
type Filter struct {
	// Program name
	Name string `yaml:"Name"`
	// Domain of program
	Domain string `yaml:"Domain"`
	// Version number
	Version string `yaml:"Version"`
	// Command to execute
	Command string `yaml:"Command"`
	// Arguments to command
	Arguments []string `yaml:"Arguments"`
	// VCS where filter can be found
	VCS `yaml:"VCS"`
}

// VCS Defines the type and location of a version control system
type VCS struct {
	// VCS type
	Type string `yaml:"Type"`
	// VCS location
	Location string `yaml:"Location"`
}

// New returns a new filter with defaults set
func New() Filter {
	return Filter{
		Arguments: make([]string, 0),
	}
}

// FilterFile returns a new Filter from a file definitions
func FilterFile(file string) (Filter, error) {
	f := New()
	return f, fromFile(file, &f)
}

// Argument appends an argument
func (f *Filter) Argument(a string) {
	f.Arguments = append(f.Arguments, a)
}

// Exec creates an executable command that is attached to input and output
func (f *Filter) Exec() *Exec {
	e := Exec{
		command: exec.Command(f.Command, f.Arguments...),
		err:     Error{command: f.Command, err: make([]byte, 0)},
	}
	e.command.Stderr = &e.err
	return &e
}

// Exec is a command and with input and output attached
type Exec struct {
	command *exec.Cmd
	link    *Exec
	writer  *io.PipeWriter
	err     Error
}

// SetInput sets input for the Exec
func (e *Exec) SetInput(r io.Reader) {
	e.command.Stdin = r
}

// SetOutput sets input for the Exec
func (e *Exec) SetOutput(w io.Writer) {
	if e.link != nil {
		e.link.SetOutput(w)
		return
	}
	e.command.Stdout = w
}

// Errors returns an array of errors
func (e *Exec) Errors() []error {
	es := make([]error, 0)
	e.errors(&es)
	return es
}

// Detach, detaches the execution from the calling process
func (e *Exec) Detach() {
	e.command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func (e *Exec) errors(es *[]error) {
	if len(e.err.err) > 0 {
		*es = append(*es, e.err)
	}
	if e.link != nil {
		e.link.errors(es)
		return
	}
}

// Run the command against the inputs and outputs
func (e *Exec) Run() error {
	err := e.command.Start()
	if err != nil {
		return err
	}
	err = e.childStart()
	if err != nil {
		return err
	}
	err = e.pipeCommands()
	if err != nil {
		return err
	}
	return nil
}

func (e *Exec) childStart() error {
	if e.link != nil {
		err := e.link.command.Start()
		if err != nil {
			return err
		}
		return e.link.childStart()
	}
	return nil
}

func (e *Exec) pipeCommands() error {
	err := e.command.Wait()
	if err != nil {
		return err
	}
	if e.link == nil {
		if e.writer != nil {
			e.writer.Close()
		}
	}
	if e.link != nil {
		if e.writer != nil {
			e.writer.Close()
		}
		err = e.link.pipeCommands()
		if err != nil {
			return err
		}

	}
	return nil

}

// Chain of commands where each feeds into stdin of the next
type Chain struct {
	Filters []Filter `yaml:"Chain"`
}

// NewChain returns a new chain
func NewChain() Chain {
	return Chain{make([]Filter, 0)}
}

// ChainFile Loads a chain from a file
func ChainFile(file string) (Chain, error) {
	c := NewChain()
	return c, fromFile(file, &c)
}

// Exec creates an exec from a chain of filters
func (c *Chain) Exec() (*Exec, error) {
	var returnExec *Exec
	var lastExec *Exec
	for k, f := range c.Filters {
		e := f.Exec()
		if k == 0 {
			returnExec = e
		} else {
			r, w := io.Pipe()
			e.command.Stdin = r
			lastExec.writer = w
			lastExec.command.Stdout = w
			lastExec.link = e
		}
		lastExec = e
	}
	return returnExec, nil
}

// Chains is a map of chains that can be retrieved by name
type Chains map[string]Chain

// NewChains returns a new map of chains
func NewChains() Chains {
	return make(Chains)
}

// ChainsFile returns Chains loaded with data from a file
func ChainsFile(file string) (Chains, error) {
	c := NewChains()
	return c, fromFile(file, &c)
}

// Get returns a chain by name or error if not found.
func (c Chains) Get(name string) (Chain, error) {
	if v, ok := c[name]; ok {
		return v, nil
	}
	return Chain{}, ErrChainDoesNotExist
}

func fromFile(file string, dest interface{}) error {
	if !isPtr(dest) {
		return ErrNotPointer
	}
	_, err := os.Stat(file)
	if err != nil {
		return err
	}
	filebytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(filebytes, dest)
	if err != nil {
		return err
	}
	return nil
}

func isPtr(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return false
	}
	return true
}
