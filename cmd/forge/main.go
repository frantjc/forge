package main

import (
	"context"
	"errors"
	"fmt"
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
		fmt.Fprintln(os.Stderr, err.Error())
	}

	stop()
	os.Exit(errorcode.ExitCode(err))
}
