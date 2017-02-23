// +build darwin

package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func FindOvftool() (string, error) {
	const name = "ovftool"
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}

	// ignore error
	home, _ := HomeDirectory()
	var defaultPaths = []string{
		"/Applications/VMware Fusion.app/Contents/Library/VMware OVF Tool/ovftool",
	}
	var vmwareDirs = []string{
		"/Applications/VMware Fusion.app",
	}
	if home != "" {
		defaultPaths = append(defaultPaths, filepath.Join(home, defaultPaths[0]))
		vmwareDirs = append(vmwareDirs, filepath.Join(home, vmwareDirs[0]))
	}
	for _, file := range defaultPaths {
		if _, err := os.Stat(file); err != nil {
			continue
		}
		if path, err := exec.LookPath(file); err == nil {
			return path, nil
		}
	}
	for _, root := range vmwareDirs {
		if _, err := os.Stat(root); err != nil {
			continue
		}
		if path, err := FindExecutable(root, name); err == nil {
			return path, nil
		}
	}
	return "", &exec.Error{Name: name, Err: exec.ErrNotFound}
}
