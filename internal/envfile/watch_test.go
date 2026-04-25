package envfile

import (
	"os"
	"testing"
	"time"
)

func writeTempWatchEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "watch_test_*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeTempWatchEnv(t, "KEY=original\n")

	done := make(chan struct{})
	defer close(done)

	ch, err := Watch(path, 50*time.Millisecond, done)
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}

	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-ch:
		if event.Err != nil {
			t.Fatalf("unexpected error in event: %v", event.Err)
		}
		if !event.Changed {
			t.Error("expected Changed=true")
		}
		if event.File != path {
			t.Errorf("expected file %q, got %q", path, event.File)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for change event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	path := writeTempWatchEnv(t, "KEY=stable\n")

	done := make(chan struct{})
	defer close(done)

	ch, err := Watch(path, 50*time.Millisecond, done)
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}

	select {
	case event := <-ch:
		t.Errorf("unexpected event received: %+v", event)
	case <-time.After(200 * time.Millisecond):
		// expected: no change
	}
}

func TestWatch_MissingFileReturnsError(t *testing.T) {
	_, err := Watch("/nonexistent/path/.env", 50*time.Millisecond, make(chan struct{}))
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWatch_StopsOnDoneSignal(t *testing.T) {
	path := writeTempWatchEnv(t, "KEY=value\n")

	done := make(chan struct{})
	ch, err := Watch(path, 30*time.Millisecond, done)
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}

	close(done)

	// Channel should close after done is signalled
	select {
	case _, open := <-ch:
		if open {
			t.Log("received event before close (acceptable)")
		}
	case <-time.After(300 * time.Millisecond):
		t.Error("channel did not close after done signal")
	}
}
