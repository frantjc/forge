package main

import (
	"os"

	"github.com/frantjc/forge/pkg/command"
	"github.com/frantjc/forge/pkg/errbubble"

	_ "gocloud.dev/blob/fileblob"
)

func main() {
	err := command.NewRoot().Execute()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}
	os.Exit(errbubble.ExitCode(err))
}
