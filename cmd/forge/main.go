package main

import (
	"context"
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

	if pErr := xerrors.Ignore(
		err,
		ore.ErrContainerExitedWithNonzeroExitCode,
	); pErr != nil {
		fmt.Fprintln(os.Stderr, pErr.Error())
	}

	stop()
	xos.ExitFromError(err)
}
