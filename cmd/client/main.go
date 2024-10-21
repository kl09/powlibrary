package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kl09/powlibrary/internal/clientapp"
	"github.com/kl09/powlibrary/internal/proto/protoconnect"
)

func main() {
	logger := slog.With("component", "main")
	logger.Info("booting up")

	client := &http.Client{}
	libClient := protoconnect.NewLibraryServiceClient(client, "http://powserver:80")

	ctx, ctxCancel := context.WithCancelCause(context.Background())

	app := clientapp.NewClientApp(
		libClient,
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
