package main

import (
	"flag"
	"log"
	"os"

	"github.com/gdbu/filelocker"
	"github.com/hatchify/errors"
)

const (
	// ErrEmptyFilename is returned when a filename is empty
	ErrEmptyFilename = errors.Error("invalid filename, cannot be empty")
)

func main() {
	var (
		filename string
		action   string
		command  string
	)

	flag.StringVar(&filename, "filename", "", "filename.txt")
	flag.StringVar(&action, "action", "lock", "action to take on file (lock, trylock)")
	flag.StringVar(&command, "command", "", "command to run under lock")
	flag.Parse()

	if len(filename) == 0 {
		log.Fatal(ErrEmptyFilename)
	}

	var (
		f   *os.File
		err error
	)

	if f, err = os.Open(filename); err != nil {
		log.Fatalf("error encountered while opening file: %v", err)
	}
	defer f.Close()
	defer filelocker.Unlock(f)

	var fn func(*os.File) error
	if fn, err = getAction(action); err != nil {
		log.Fatal(err)
	}

	if err = fn(f); err != nil {
		log.Fatal(err)
	}

	if err = runCommand(command); err != nil {
		log.Fatal(err)
	}
}
