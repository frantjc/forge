package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/forge/command"
	"github.com/frantjc/forge/ore"
	errorcode "github.com/frantjc/go-error-code"
)

func main() {
	var (
		ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		err       error
	)

	if err = command.NewForge().ExecuteContext(ctx); err != nil && !errors.Is(err, ore.ErrContainerExitedWithNonzeroExitCode) {
		os.Stderr.WriteString(err.Error() + "\n")
	}

	stop()
	os.Exit(errorcode.ExitCode(err))
}
