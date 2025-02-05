package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/command"
	xos "github.com/frantjc/x/os"
)

func main() {
	var (
		cmd       = command.NewForge()
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	)

	cmd.Version = SemVer()

	err := cmd.ExecuteContext(ctx)
	if err != nil && !errors.Is(err, forge.ErrContainerExitedWithNonzeroExitCode) {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	stop()
	xos.ExitFromError(err)
}
