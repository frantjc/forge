package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/forge/command"
	errorcode "github.com/frantjc/go-error-code"

	_ "gocloud.dev/blob/fileblob"
)

func main() {
	var (
		ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		err       error
	)

	if err = command.New().ExecuteContext(ctx); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}

	stop()
	os.Exit(errorcode.ExitCode(err))
}
