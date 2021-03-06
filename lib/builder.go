package gin

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type Builder interface {
	Build() error
	Binary() string
	Errors() string
}

type builder struct {
	dir       string
	binary    string
	errors    string
	useGodep  bool
	buildArgs []string
	mutex     *sync.Mutex
}

func NewBuilder(dir string, bin string, useGodep bool, buildArgs []string) Builder {
	if len(bin) == 0 {
		bin = "bin"
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") { // check if it already has the .exe extension
			bin += ".exe"
		}
	}

	m := &sync.Mutex{}

	return &builder{dir: dir, binary: bin, useGodep: useGodep, buildArgs: buildArgs, mutex: m}
}

func (b *builder) Binary() string {
	return b.binary
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	args := append([]string{"go", "build", "-o", b.binary}, b.buildArgs...)

	var command *exec.Cmd
	if b.useGodep {
		args = append([]string{"godep"}, args...)
	}
	b.mutex.Lock()
	command = exec.Command(args[0], args[1:]...)
	b.mutex.Unlock()

	command.Dir = b.dir

	output, err := command.CombinedOutput()

	if command.ProcessState.Success() {
		b.errors = ""
	} else {
		b.errors = string(output)
	}

	if len(b.errors) > 0 {
		return fmt.Errorf(b.errors)
	}

	return err
}
