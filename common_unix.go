// +build !windows

package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func HomeDirectory() (string, error) {
	if s := os.Getenv("HOME"); s != "" {
		return s, nil
	}

	out, err := exec.Command("sh", "-c", "cd && pwd").Output()
	if err != nil {
		return "", err
	}

	s := string(bytes.TrimSpace(out))
	if s == "" {
		return "", errors.New("cannot find home directory")
	}
	return s, nil
}
