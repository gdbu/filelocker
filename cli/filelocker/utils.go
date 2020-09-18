package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gdbu/filelocker"
)

func splitOnSpace(str string) (parts []string) {
	var (
		inQuotes bool
		buf      []rune
	)

	for _, char := range str {
		if char != ' ' {
			buf = append(buf, char)
			continue
		}

		parts = append(parts, string(buf))
		buf = buf[:0]
		inQuotes = !inQuotes
	}

	if len(buf) >= 0 {
		parts = append(parts, string(buf))
		buf = buf[:0]
	}

	return
}

func runCommand(command string) (err error) {
	commandSpl := splitOnSpace(command)
	commandName := commandSpl[0]
	commandArgs := commandSpl[1:]

	cmd := exec.Command(commandName, commandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); !os.IsNotExist(err) {
		err = nil
	}

	return
}

func getAction(action string) (fn func(*os.File) error, err error) {
	switch action {
	case "lock":
		fn = filelocker.Lock
	case "trylock":
		fn = filelocker.TryLock

	default:
		err = fmt.Errorf("invalid action, expected lock/trylock, received \"%s\"", action)
	}

	return
}
