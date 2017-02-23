package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func workstationInstallPaths() ([]string, error) {
	const keypath = `SOFTWARE\Wow6432Node\VMware, Inc.\VMware Workstation`
	const keymode = registry.READ | registry.ENUMERATE_SUB_KEYS
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keypath, keymode)
	if err != nil {
		return nil, err
	}
	var paths []string
	if s, _, err := key.GetStringValue("InstallPath64"); err == nil {
		paths = append(paths, s)
	}
	if s, _, err := key.GetStringValue("InstallPath"); err == nil {
		if len(paths) == 0 || paths[0] != s {
			paths = append(paths, s)
		}
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("reading registry key (%s): %s", keypath, err)
	}
	return paths, nil
}

func FindOvftool() (string, error) {
	const name = "ovftool.exe"
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}

	installPaths, err := workstationInstallPaths()
	if err != nil {
		return "", err
	}
	for _, dir := range installPaths {
		file := filepath.Join(dir, "ovftool", "ovftool.exe")
		if path, err := exec.LookPath(file); err == nil {
			return path, nil
		}
		if path, err := FindExecutable(dir, name); err == nil {
			return path, nil
		}
	}
	return "", &exec.Error{Name: name, Err: exec.ErrNotFound}
}
