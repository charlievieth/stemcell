package main

import "testing"

func TestReadVMXDir(t *testing.T) {
	var tests = []struct {
		dirname string
		version int
		expErr  bool
	}{
		{"valid", 11, false},
	}
	_ = tests

	// for _, x := range tests {
	// }
}
