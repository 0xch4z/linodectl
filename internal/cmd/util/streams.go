package util

import "io"

type IOStreams struct {
	In       io.Reader
	Out, Err io.Writer
}
