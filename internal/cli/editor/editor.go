package editor

// Inspired by: https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/util/editor/editor.go

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/strutil"
	"github.com/moby/term"
)

func NewDefaultEditor() Editor {
	editor := strutil.Fallback(os.Getenv("EDITOR"), defaultEditor)
	shell := strutil.Fallback(os.Getenv("SHELL"), defaultShell)
	args := append(defaultShellArgs, editor)
	return Editor{Args: append([]string{shell}, args...)}
}

type Editor struct {
	Args []string
}

func (e Editor) args(path string) []string {
	args := make([]string, len(e.Args))
	copy(args, e.Args)
	last := args[len(args)-1]
	args[len(args)-1] = fmt.Sprintf("%s %q", last, path)
	return args
}

func (Editor) backupInBuffer(in io.Reader) func() {
	inFD, isTerminal := term.GetFdInfo(in)
	if !isTerminal {
		if f, err := os.Open("/dev/tty"); err == nil {
			defer f.Close()
			inFD = f.Fd()
			isTerminal = term.IsTerminal(inFD)
		}
	}

	if !isTerminal {
		return nil
	}

	state, err := term.SaveState(inFD)
	if err != nil {
		return nil
	}

	return func() {
		_ = term.RestoreTerminal(inFD, state)
	}
}

func (e Editor) Launch(ioStreams util.IOStreams, path string) error {
	if len(e.Args) == 0 {
		panic("no editor defined")
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	restore := e.backupInBuffer(ioStreams.In)
	defer func() {
		if restore != nil {
			restore()
		}
	}()

	args := e.args(abs)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = ioStreams.Out
	cmd.Stderr = ioStreams.Err
	cmd.Stdin = ioStreams.In
	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.Error); ok {
			if err.Err == exec.ErrNotFound {
				panic(err)
			}
		}
		panic(err)
	}
	return nil
}

func (e Editor) EditReader(prefix, suffix string, ioStreams util.IOStreams, r io.Reader) ([]byte, string, error) {
	// create temporary file for editing
	f, err := os.CreateTemp("", prefix+"*"+suffix)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	// copy the contents of io.Reader into the temporary file
	path := f.Name()
	if _, err := io.Copy(f, r); err != nil {
		// clean up the file
		os.Remove(path)
		return nil, path, err
	}

	f.Close()

	// open the temporary file with the user's editor
	if err := e.Launch(ioStreams, path); err != nil {
		return nil, path, err
	}
	bytes, err := ioutil.ReadFile(path)
	return bytes, path, err
}
