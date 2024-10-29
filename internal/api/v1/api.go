package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kl09/powlibrary/internal/domain"
	"golang.org/x/time/rate"
)

type QuotesHandler struct {
	libService   domain.LibraryService
	powService   domain.POWService
	tasksStorage domain.TasksStorage
	rateLimiter  *rate.Limiter

	logger *slog.Logger
}

func NewQuotesHandler(
	libService domain.LibraryService,
	powService domain.POWService,
	tasksStorage domain.TasksStorage,
	maxRPS int,
	logger *slog.Logger,
) *QuotesHandler {
	return &QuotesHandler{
		libService:   libService,
		powService:   powService,
		tasksStorage: tasksStorage,
		rateLimiter:  rate.NewLimiter(rate.Limit(maxRPS), maxRPS),
		logger:       logger,
	}
}

func (h *QuotesHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !h.rateLimiter.Allow() {
		h.logger.Info("rate limit exceeded")
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	h.logger.Info(fmt.Sprintf("request: %s", r.URL.Path))
	switch r.URL.Path {
	case "/GenerateTask":
		h.GenerateTask(w, r)
	case "/GetQuote":
		h.GetQuote(w, r)
	default:
		http.NotFound(w, r)
	}
}

// GenerateTask generates a task for PoW.
func (h *QuotesHandler) GenerateTask(w http.ResponseWriter, r *http.Request) {
	var jsonReq struct {
		UserID string `json:"user_id"`
	}

	// read 1MB max from the body - to prevent DoS attacks.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	err := json.NewDecoder(r.Body).Decode(&jsonReq)
	if err != nil {
		returnErr(w, http.StatusInternalServerError, err)
		return
	}

	if jsonReq.UserID == "" {
		returnErr(w, http.StatusBadRequest, errors.New("user_id is required"))
		return
	}

	task := h.tasksStorage.GetForUser(jsonReq.UserID)
	if task == nil || task.ResolvedAt != nil {
		taskCode, difficulty, err := h.powService.Generate()
		if err != nil {
			returnErr(w, http.StatusInternalServerError, err)
			return
		}

		task = &domain.POWTask{
			Task:       taskCode,
			UserID:     jsonReq.UserID,
			Difficulty: difficulty,
		}
		h.tasksStorage.Add(*task)
	}

	type resp struct {
		Task       string `json:"task"`
		Difficulty int    `json:"difficulty"`
	}

	response := resp{
		Task:       task.Task,
		Difficulty: task.Difficulty,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		returnErr(w, http.StatusInternalServerError, err)
		return
	}
}

// GetQuote returns a random quote if PoW is validated.
func (h *QuotesHandler) GetQuote(w http.ResponseWriter, r *http.Request) {
	var jsonReq struct {
		UserID string `json:"user_id"`
		Task   string `json:"task"`
		Hash   string `json:"hash"`
	}

	// read 1MB max from the body - to prevent DoS attacks.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	err := json.NewDecoder(r.Body).Decode(&jsonReq)
	if err != nil {
		returnErr(w, http.StatusInternalServerError, err)
		return
	}

	if jsonReq.UserID == "" {
		returnErr(w, http.StatusBadRequest, errors.New("user_id is required"))
		return
	}
	if jsonReq.Task == "" {
		returnErr(w, http.StatusBadRequest, errors.New("tash is required"))
		return
	}
	if jsonReq.Hash == "" {
		returnErr(w, http.StatusBadRequest, errors.New("hash is required"))
		return
	}

	task := h.tasksStorage.GetForUser(jsonReq.UserID)
	if task == nil || task.ResolvedAt != nil {
		returnErr(w, http.StatusInternalServerError, errors.New("task isn't generated yet"))
		return
	}

	if task.Task != jsonReq.Task {
		returnErr(w, http.StatusBadRequest, errors.New("task is invalid"))
		return
	}

	err = h.powService.Validate(jsonReq.Task, jsonReq.Hash)
	if err != nil {
		returnErr(w, http.StatusInternalServerError, err)
		return
	}

	h.tasksStorage.MarkAsUsed(*task)

	type resp struct {
		Quote string `json:"quote"`
	}

	response := resp{
		Quote: h.libService.GetRandomQuote(r.Context()),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		returnErr(w, http.StatusInternalServerError, err)
		return
	}
}

func returnErr(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(fmt.Sprintf("err: %s", err.Error())))
}

// WrapWithTimeoutHandler wraps handler `h` so that it runs with the given time limit `dt`.
func WrapWithTimeoutHandler(h http.Handler, dt time.Duration) http.Handler {
	msg := fmt.Sprintf(
		`{"error":{"code":"%s","message":"%s"}}`,
		"request_timeout",
		"Request timed out.",
	)
	return http.TimeoutHandler(h, dt, msg)
}
