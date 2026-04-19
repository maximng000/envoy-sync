package cmd

import (
	"bytes"
	"io"
	"os"
)

// captureOutput runs fn and returns everything written to stdout.
// It temporarily replaces os.Stdout with a pipe to capture output.
func captureOutput(fn func() error) (string, error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	fnErr := fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), fnErr
}
