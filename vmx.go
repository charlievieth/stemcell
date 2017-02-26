package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// A VMX defines a versioned .vmx file
type VMX struct {
	Dirname string // absolute directory path
	Path    string // absolute file path
	Name    string // file name
	Version int    // file version
}

type byVersion []VMX

func (v byVersion) Len() int           { return len(v) }
func (v byVersion) Less(i, j int) bool { return v[i].Version < v[j].Version }
func (v byVersion) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

func findVMXFile(dirname string) (string, error) {
	pattern := filepath.Join(dirname, "*.vmx")
	names, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("vmx file not found in directory: %s", dirname)
	}
	if len(names) != 1 {
		return "", fmt.Errorf("multiple vmx files (%v) found in directory: %s", names, dirname)
	}
	name := names[0]
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	if fi.IsDir() {
		return "", fmt.Errorf("directory (%s) contains a .vmx directory (%s) expeted a file",
			dirname, name)
	}
	return name, nil
}

// ReadVMXDir returns the VMX entries in directory dirname that match the string
// format.  The format is used to scan the directory version (ex: "vmx-%d") and
// match files ("vmx-%d" is converted to "vmx-*").
//
// The returned VMX entries are sorted by version.
func ReadVMXDir(dirname, format string) ([]VMX, error) {
	var err error
	dirname, err = filepath.Abs(dirname)
	if err != nil {
		return nil, err
	}
	if strings.Count(format, "%d") != 1 {
		return nil, fmt.Errorf("invalid format string: %s", format)
	}
	pattern := strings.Replace(format, "%d", "*", -1)
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	if len(fis) == 0 {
		return nil, fmt.Errorf("empty vmx directory: %s", dirname)
	}
	var files []VMX
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		name := fi.Name()
		match, err := filepath.Match(pattern, name)
		if err != nil {
			// malformed pattern - likely due to the 'format' argument
			return nil, err
		}
		if !match {
			continue
		}
		var v int
		n, err := fmt.Sscanf(name, format, &v)
		if err != nil {
			continue // TODO (CEV): return error
		}
		if n != 1 {
			continue // TODO (CEV): return error
		}
		dir := filepath.Join(dirname, name)
		path, err := findVMXFile(dir)
		if err != nil {
			continue // TODO (CEV): log error
		}
		files = append(files, VMX{
			Dirname: dir,
			Path:    path,
			Name:    filepath.Base(path),
			Version: v,
		})
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("empty vmx directory: %s", dirname)
	}
	sort.Sort(byVersion(files))
	return files, nil
}
