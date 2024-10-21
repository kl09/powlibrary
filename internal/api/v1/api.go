package v1

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/kl09/powlibrary/internal/domain"
	v1 "github.com/kl09/powlibrary/internal/proto"
)

type QuotesHandler struct {
	libService   domain.LibraryService
	powService   domain.POWService
	tasksStorage domain.TasksStorage
}

func NewQuotesHandler(
	libService domain.LibraryService,
	powService domain.POWService,
	tasksStorage domain.TasksStorage,
) *QuotesHandler {
	return &QuotesHandler{
		libService:   libService,
		powService:   powService,
		tasksStorage: tasksStorage,
	}
}

// GenerateTask generates a task for PoW.
func (h *QuotesHandler) GenerateTask(
	ctx context.Context, req *connect.Request[v1.GenerateTaskRequest]) (
	*connect.Response[v1.GenerateTaskResponse], error,
) {
	if req.Msg.GetUserId() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user_id is required"))
	}

	task := h.tasksStorage.GetForUser(req.Msg.GetUserId())
	if task == nil || task.IsUsed {
		taskCode, difficulty, err := h.powService.Generate()
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		task = &domain.POWTask{
			Task:       taskCode,
			UserID:     req.Msg.GetUserId(),
			Difficulty: difficulty,
		}
		h.tasksStorage.Add(*task)
	}

	return connect.NewResponse(&v1.GenerateTaskResponse{
		Task:       task.Task,
		Difficulty: int32(task.Difficulty),
	}), nil
}

// GetQuote returns a random quote if PoW is validated.
func (h *QuotesHandler) GetQuote(
	ctx context.Context, req *connect.Request[v1.GetQuoteRequest],
) (*connect.Response[v1.GetQuoteResponse], error) {
	if req.Msg.GetUserId() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user_id is required"))
	}
	if req.Msg.GetTask() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("task is required"))
	}
	if req.Msg.GetHash() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("hash is required"))
	}

	task := h.tasksStorage.GetForUser(req.Msg.GetUserId())
	if task == nil || task.IsUsed {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("task is not generated"))
	}

	if task.Task != req.Msg.GetTask() {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("task is invalid"))
	}

	err := h.powService.Validate(req.Msg.GetTask(), req.Msg.GetHash())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	h.tasksStorage.MarkAsUsed(*task)

	return connect.NewResponse(&v1.GetQuoteResponse{
		Quote: h.libService.GetRandomQuote(ctx),
	}), nil
}
