package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kl09/powlibrary/internal/clientapp"
)

func main() {
	logger := slog.With("component", "main")
	logger.Info("booting up")
	ctx, ctxCancel := context.WithCancelCause(context.Background())

	app := clientapp.NewClientApp(
		"http://powserver:80",
		logger,
	)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		fmt.Fprint(os.Stderr, "signal received - terminating\n")
		signal.Reset()
		ctxCancel(errors.New("signal received"))
	}()

	app.Run(ctx)

	os.Exit(0)
}
