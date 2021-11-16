//go:build !windows
// +build !windows

package editor

const (
	defaultEditor = "vi"
	defaultShell  = "/bin/bash"
)

var defaultShellArgs = []string{"-c"}
