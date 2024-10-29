package clientapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kl09/powlibrary/internal/utils"
)

type ClientApp struct {
	powServerURL string
	logger       *slog.Logger
}

func NewClientApp(
	powServerURL string,
	logger *slog.Logger,
) *ClientApp {
	return &ClientApp{
		powServerURL: powServerURL,
		logger:       logger,
	}
}

func (a *ClientApp) Run(ctx context.Context) {
	a.logger.Info("client: running")

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
				}

				userID := fmt.Sprintf("test-%s", uuid.New())

				startedAt := time.Now()

				task, difficulty, err := a.generateTask(ctx, userID)
				if err != nil {
					a.logger.Error(fmt.Sprintf("client: failed to generate task: %s", err))
					continue
				}

				// ofc, in real project it won't be called from utils.
				hash, err := utils.GeneratePOW(ctx, task, difficulty)
				if err != nil {
					a.logger.Error(fmt.Sprintf("client: failed to generate PoW: %s", err))
					continue
				}

				quote, err := a.getQuote(ctx, task, hash, userID)
				if err != nil {
					a.logger.Error(fmt.Sprintf("client: failed to get quote: %s", err))
					continue
				}
				a.logger.Info(fmt.Sprintf("got quote with hash: %s, task: %s and difficulty: %d for user_id: %s", hash, task, difficulty, userID))
				a.logger.Info(fmt.Sprintf("client: quote received: %s for %f seconds", quote, time.Since(startedAt).Seconds()))
			}
		}()
	}

	wg.Wait()
}

func (a *ClientApp) generateTask(ctx context.Context, userID string) (task string, difficulty int, err error) {
	req, err := http.NewRequestWithContext(
		ctx, "POST", a.powServerURL+"/GenerateTask",
		strings.NewReader(fmt.Sprintf(`{"user_id": "%s"}`, userID)),
	)
	if err != nil {
		return "", 0, fmt.Errorf("failed to do request: %w", err)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("do request: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return "", 0, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	var resp struct {
		Task       string `json:"task"`
		Difficulty int    `json:"difficulty"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return resp.Task, resp.Difficulty, nil
}

func (a *ClientApp) getQuote(ctx context.Context, task, hash, userID string) (quote string, err error) {
	var jsonReq struct {
		UserID string `json:"user_id"`
		Task   string `json:"task"`
		Hash   string `json:"hash"`
	}
	jsonReq.Task = task
	jsonReq.Hash = hash
	jsonReq.UserID = userID

	b, err := json.Marshal(jsonReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx, "POST", a.powServerURL+"/GetQuote", bytes.NewReader(b),
	)
	if err != nil {
		return "", fmt.Errorf("failed to do request: %w", err)
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	var resp struct {
		Quote string `json:"quote"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return resp.Quote, nil
}
