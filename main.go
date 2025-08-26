package main

import (
	"context"
	"os"
	"syscall"

	"github.com/riabininkf/go-modules/cmd"

	_ "github.com/riabininkf/http-auth-example/cmd"
)

func main() {
	ctx, cancel := cmd.ContextFromSignals(
		context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)

	defer cancel()

	if err := cmd.Execute(ctx); err != nil {
		os.Exit(1)
	}
}
