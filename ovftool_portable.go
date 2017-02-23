// +build !darwin,!windows

package main

import "os/exec"

func FindOvftool() (string, error) {
	const name = "ovftool"
	return exec.LookPath(name)
}
