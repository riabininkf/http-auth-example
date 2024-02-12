package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/riabininkf/go-project-template/cmd"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		cancelFunc()
	}()

	if err := cmd.RootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
