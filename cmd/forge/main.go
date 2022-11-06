package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/forge/pkg/command"
	"github.com/frantjc/forge/pkg/errbubble"

	_ "gocloud.dev/blob/fileblob"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		sigC        = make(chan os.Signal, 1)
		err         error
	)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigC
		cancel()
	}()

	if err = command.New().ExecuteContext(ctx); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}

	os.Exit(errbubble.ExitCode(err))
}
