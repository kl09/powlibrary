package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	v1 "github.com/kl09/powlibrary/internal/api/v1"
	"github.com/kl09/powlibrary/internal/database"
	"github.com/kl09/powlibrary/internal/library"
	"github.com/kl09/powlibrary/internal/pow"
)

const (
	defaultPOWDifficulty = 3
	maxRPS               = 100
)

func main() {
	logger := slog.With("component", "main")
	logger.Info("booting up")

	tasksStorage := database.NewTasksStorage(time.Minute)
	powService := pow.NewProofOfWork(defaultPOWDifficulty)

	publicHandler := v1.NewQuotesHandler(
		library.NewLibrary(), powService, tasksStorage, maxRPS, logger,
	)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(publicHandler.Handler))

	server := &http.Server{
		Addr:              ":80",
		Handler:           v1.WrapWithTimeoutHandler(mux, 5*time.Second),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	defer server.Close()

	ctx, ctxCancel := context.WithCancelCause(context.Background())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		fmt.Fprint(os.Stderr, "signal received - terminating\n")
		signal.Reset()
		ctxCancel(errors.New("signal received"))
		server.Close()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(6 * time.Second):
				avgTime := tasksStorage.AvgTimeToResolve()
				logger.Info(fmt.Sprintf("server: avg time to resolve: %f", avgTime))
				if avgTime > 5 {
					cur := powService.DecreaseDifficulty()
					logger.Info(fmt.Sprintf("difficulty decreased to %d", cur))
				} else if avgTime < 2 {
					cur := powService.IncreaseDifficulty()
					logger.Info(fmt.Sprintf("difficulty increased to %d", cur))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(5 * time.Second):
				tasksStorage.ClearCache()
			case <-ctx.Done():
				return
			}
		}
	}()

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
