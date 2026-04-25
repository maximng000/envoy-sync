package envfile

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent describes a change detected in a watched .env file.
type WatchEvent struct {
	File    string
	Changed bool
	Err     error
}

// fileChecksum returns the MD5 checksum of a file's contents.
func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Watch polls the given file at the specified interval, sending a WatchEvent
// on the returned channel whenever the file changes or an error occurs.
// The caller must close the done channel to stop watching.
func Watch(path string, interval time.Duration, done <-chan struct{}) (<-chan WatchEvent, error) {
	initial, err := fileChecksum(path)
	if err != nil {
		return nil, fmt.Errorf("watch: initial checksum failed: %w", err)
	}

	ch := make(chan WatchEvent, 1)

	go func() {
		defer close(ch)
		last := initial
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				current, err := fileChecksum(path)
				if err != nil {
					ch <- WatchEvent{File: path, Err: err}
					continue
				}
				if current != last {
					last = current
					ch <- WatchEvent{File: path, Changed: true}
				}
			}
		}
	}()

	return ch, nil
}
