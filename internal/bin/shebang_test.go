package bin_test

import (
	"fmt"
	"testing"

	"github.com/frantjc/forge/internal/bin"
)

func TestHasShebang(t *testing.T) {
	var (
		script = `#!/bin/bash

echo hello
`
		expected = true
		actual   = bin.HasShebang(script)
	)

	if expected != actual {
		t.Error("actual value", `"`+fmt.Sprint(actual)+`"`, "does not equal expected value", `"`+fmt.Sprint(expected)+`"`)
		t.FailNow()
	}
}

func TestDoesNotHaveShebang(t *testing.T) {
	var (
		script   = "echo hello"
		expected = false
		actual   = bin.HasShebang(script)
	)

	if expected != actual {
		t.Error("actual value", `"`+fmt.Sprint(actual)+`"`, "does not equal expected value", `"`+fmt.Sprint(expected)+`"`)
		t.FailNow()
	}
}
