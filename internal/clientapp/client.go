package clientapp

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	v1 "github.com/kl09/powlibrary/internal/proto"
	"github.com/kl09/powlibrary/internal/proto/protoconnect"
	"github.com/kl09/powlibrary/internal/utils"
)

type ClientApp struct {
	client protoconnect.LibraryServiceClient
	logger *slog.Logger
}

func NewClientApp(
	client protoconnect.LibraryServiceClient,
	logger *slog.Logger,
) *ClientApp {
	return &ClientApp{
		client: client,
		logger: logger,
	}
}

func (a *ClientApp) Run(ctx context.Context) {
	a.logger.Info("client: running")
	var userID = "test"

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		resp, err := a.client.GenerateTask(ctx, connect.NewRequest(&v1.GenerateTaskRequest{
			UserId: userID,
		}))
		if err != nil {
			a.logger.Error(fmt.Sprintf("client: failed to generate task: %s", err))
			continue
		}

		// ofc, in real project it won't be called from utils.
		hash, err := utils.GeneratePOW(ctx, resp.Msg.GetTask(), int(resp.Msg.GetDifficulty()))
		if err != nil {
			a.logger.Error(fmt.Sprintf("client: failed to generate PoW: %s", err))
			continue
		}

		respQ, err := a.client.GetQuote(ctx, connect.NewRequest(&v1.GetQuoteRequest{
			Task:   resp.Msg.GetTask(),
			Hash:   hash,
			UserId: userID,
		}))
		if err != nil {
			a.logger.Error(fmt.Sprintf("client: failed to get quote: %s", err))
			continue
		}
		a.logger.Info(fmt.Sprintf("got quote with hash: %s, task: %s and difficulty: %d", hash, resp.Msg.GetTask(), resp.Msg.Difficulty))
		a.logger.Info(fmt.Sprintf("client: quote received: %s", respQ.Msg.GetQuote()))
	}
}
