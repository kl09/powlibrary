package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"connectrpc.com/connect"
	"github.com/kl09/powlibrary/internal/api/middleware"
	v1 "github.com/kl09/powlibrary/internal/api/v1"
	"github.com/kl09/powlibrary/internal/database"
	"github.com/kl09/powlibrary/internal/library"
	"github.com/kl09/powlibrary/internal/pow"
	"github.com/kl09/powlibrary/internal/proto/protoconnect"
)

const powDifficulty = 5

func main() {
	logger := slog.With("component", "main")
	logger.Info("booting up")

	publicHandler := v1.NewQuotesHandler(
		library.NewLibrary(), pow.NewProofOfWork(powDifficulty), database.NewTasksStorage(),
	)
	route, handler := protoconnect.NewLibraryServiceHandler(
		publicHandler,
		connect.WithRecover(middleware.Recover),
		connect.WithInterceptors(
			middleware.Logging(middleware.WithLogger(slog.With("component", "api"))),
		),
	)

	mux := http.NewServeMux()
	mux.Handle(route, handler)

	server := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}
	defer server.Close()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		fmt.Fprint(os.Stderr, "signal received - terminating\n")
		signal.Reset()
		server.Close()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info(fmt.Sprintf("server: listening to %s", server.Addr))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(fmt.Sprintf("server: serve finished with err: %s", err))
		}
		logger.Info("server: stopped")
	}()

	wg.Wait()

	os.Exit(0)
}
