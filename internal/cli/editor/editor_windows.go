//go:build unix
// +build unix

package editor

const (
	defaultEditor = "notepad"
	defaultShell  = "cmd"
)

var defaultShellArgs = []string{}
