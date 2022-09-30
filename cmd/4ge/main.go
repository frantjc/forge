package main

import (
	"os"

	"github.com/frantjc/forge/pkg/command"
	"github.com/frantjc/forge/pkg/errbubble"

	_ "gocloud.dev/blob/fileblob"
)

func main() {
	os.Exit(errbubble.ExitCode(command.NewRoot().Execute()))
}
