package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	// set := flag.NewFlagSet("base", flag.ContinueOnError)
	// _ = set

	fs, err := ReadVMXDir("testdata/vmx_tests/valid", "vmx-%d")
	if err != nil {
		Fatal(err)
	}
	PrintJSON(fs)
	return 0
}

// type Config struct {
// ADMINISTRATOR_PASSWORD
// BUILDER_PATH
// OUTPUT_DIR
// VERSION
// VMX_DIR
// }

func PrintJSON(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "    ")
	if err == nil {
		fmt.Println(string(b))
	}
	return err
}

func Fatal(err interface{}) {
	if err == nil {
		return
	}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		file = filepath.Base(file)
	}
	switch err.(type) {
	case error, string, fmt.Stringer:
		if ok {
			fmt.Fprintf(os.Stderr, "Error (%s:%d): %s", file, line, err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
		}
	default:
		if ok {
			fmt.Fprintf(os.Stderr, "Error (%s:%d): %#v\n", file, line, err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %#v\n", err)
		}
	}
	os.Exit(1)
}
