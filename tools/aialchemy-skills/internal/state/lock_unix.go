//go:build unix

package state

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// Lock is an exclusive advisory process lock held for the duration of every
// mutating command.
type Lock struct {
	f    *os.File
	path string
}

// LockHeldError reports that another process holds the lock. Because flock
// locks are released by the kernel when the holder dies, a held lock always
// belongs to an apparently live process; it is reported, never deleted.
type LockHeldError struct {
	Path   string
	Holder string
}

func (e *LockHeldError) Error() string {
	return fmt.Sprintf("another aialchemy-skills process (%s) holds the lock %s; the kernel releases it automatically if that process has exited, so a persistent failure means a mutating command is still running", e.Holder, e.Path)
}

// AcquireLock takes the exclusive state lock, or returns *LockHeldError.
func AcquireLock(stateDir string) (*Lock, error) {
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return nil, fmt.Errorf("cannot create state directory: %w", err)
	}
	fi, err := os.Lstat(stateDir)
	if err != nil {
		return nil, fmt.Errorf("cannot inspect state directory: %w", err)
	}
	if !fi.IsDir() || fi.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("state directory %s is not a real directory", stateDir)
	}
	path := filepath.Join(stateDir, "lock")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|syscall.O_NOFOLLOW, 0o644)
	if err != nil {
		return nil, fmt.Errorf("cannot open lock file: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		content, _ := io.ReadAll(io.LimitReader(f, 256))
		f.Close()
		holder := strings.TrimSpace(string(content))
		if holder == "" {
			holder = "unknown"
		}
		return nil, &LockHeldError{Path: path, Holder: holder}
	}
	// The note is advisory, for diagnostics only; flock is the actual lock.
	f.Truncate(0)
	f.Seek(0, io.SeekStart)
	fmt.Fprintf(f, "pid %d\n", os.Getpid())
	f.Sync()
	return &Lock{f: f, path: path}, nil
}

// Release drops the lock. The lock file itself is left in place to avoid
// unlink races with concurrent acquirers.
func (l *Lock) Release() error {
	if l == nil || l.f == nil {
		return nil
	}
	syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
	err := l.f.Close()
	l.f = nil
	return err
}
