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
	xerrors "github.com/frantjc/x/errors"
	xos "github.com/frantjc/x/os"
)

func main() {
	var (
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		err       = xerrors.Ignore(
			command.NewForge().ExecuteContext(ctx),
			context.Canceled,
		)
	)

	if err != nil && !errors.Is(err, ore.ErrContainerExitedWithNonzeroExitCode) {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	stop()
	xos.ExitFromError(err)
}
