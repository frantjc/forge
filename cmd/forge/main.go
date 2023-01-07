package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/forge/command"
	errorcode "github.com/frantjc/go-error-code"
)

func main() {
	var (
		ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		err       error
	)

	if err = command.NewForge().ExecuteContext(ctx); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}

	stop()
	os.Exit(errorcode.ExitCode(err))
}
