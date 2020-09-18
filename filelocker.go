package filelocker

import (
	"context"
	"os"
	"syscall"

	"github.com/hatchify/errors"
)

const (
	// ErrIsLocked is returned when a file is already locked
	ErrIsLocked = errors.Error("file is already locked")
	// ErrTimeout is returned when a timeout has been exceeded while acquiring a lock
	ErrTimeout = errors.Error("timeout exceeded when acquiring lock")
)

// Lock will acquire a *nix file lock on the provided file
func Lock(f *os.File) (err error) {
	fd := int(f.Fd())
	return syscall.Flock(fd, syscall.LOCK_EX)
}

// TryLock will attempt to immediately acquire a *nix file lock on the provided file
func TryLock(f *os.File) (err error) {
	fd := int(f.Fd())
	if err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB); err == syscall.EWOULDBLOCK {
		return ErrIsLocked
	}

	return
}

// LockWithContext will acquire a *nix file lock on the provided file with context.Context
func LockWithContext(ctx context.Context, f *os.File) error {
	result := make(chan error)
	go func() {
		result <- Lock(f)
	}()

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		go Unlock(f)
		return ErrTimeout
	}
}

// Unlock will release a *nix file lock on the provided file
func Unlock(f *os.File) (err error) {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
