package cmd

import (
	"bytes"
	"io"
	"os"
)

// captureOutput runs fn and returns everything written to stdout.
func captureOutput(fn func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}
