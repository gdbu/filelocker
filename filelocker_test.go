package filelocker

import (
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	var (
		f   *os.File
		err error
	)

	if f, err = os.Create("test_file"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test_file")
	defer f.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	errC := make(chan error)
	doneC := make(chan struct{})
	endTime := time.Now().Add(time.Second * 3)

	go testCommand(&wg, errC, "sleep 3")
	time.Sleep(time.Millisecond)
	go testLocal(&wg, errC, f, Lock)

	go func() {
		wg.Wait()
		doneC <- struct{}{}
	}()

	select {
	case <-doneC:
	case err = <-errC:
		t.Fatal(err)
	}

	now := time.Now()
	if now.Before(endTime) {
		t.Fatalf("test ended sooner than expected by %v", endTime.Sub(now))
	}
}

func TestTryLock_with_contention(t *testing.T) {
	var (
		f   *os.File
		err error
	)

	if f, err = os.Create("test_file"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test_file")
	defer f.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	errC := make(chan error)
	doneC := make(chan struct{})

	go testCommand(&wg, errC, "sleep 3")
	time.Sleep(time.Millisecond)
	go testLocal(&wg, errC, f, TryLock)

	go func() {
		wg.Wait()
		doneC <- struct{}{}
	}()

	select {
	case <-doneC:
	case err = <-errC:
	}

	if err != ErrTimeout {
		t.Fatalf("invalid error, expected %v and received %v", ErrTimeout, err)
	}
}

func testLocal(wg *sync.WaitGroup, errC chan error, f *os.File, fn func(*os.File) error) {
	if err := fn(f); err != nil {
		errC <- err
	}

	if err := Unlock(f); err != nil {
		errC <- err
	}

	wg.Done()
}

func testCommand(wg *sync.WaitGroup, errC chan error, command string) {
	cmd := exec.Command("./cli/filelocker/filelocker", "--filename", "./test_file", "--command", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		errC <- err
	}

	wg.Done()
}
