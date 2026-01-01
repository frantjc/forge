package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/forge/command"
	xos "github.com/frantjc/x/os"
)

func main() {
	var (
		cmd       = command.NewShim()
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	)

	err := cmd.ExecuteContext(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	stop()
	xos.ExitFromError(err)
}
