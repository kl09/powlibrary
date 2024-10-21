package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/kl09/powlibrary/internal/database"
	"github.com/kl09/powlibrary/internal/library"
	"github.com/kl09/powlibrary/internal/pow"
	v1 "github.com/kl09/powlibrary/internal/proto"
	"github.com/kl09/powlibrary/internal/proto/protoconnect"
	"github.com/stretchr/testify/require"
)

func TestQuotesHandler_GenerateTask(t *testing.T) {
	route, handler := protoconnect.NewLibraryServiceHandler(
		NewQuotesHandler(library.NewLibrary(), pow.NewProofOfWork(1), database.NewTasksStorage()),
	)

	mux := http.NewServeMux()
	mux.Handle(route, handler)

	s := httptest.NewServer(mux)
	t.Cleanup(s.Close)

	client := protoconnect.NewLibraryServiceClient(s.Client(), s.URL)

	t.Run("success", func(t *testing.T) {
		resp, err := client.GenerateTask(context.Background(), connect.NewRequest(&v1.GenerateTaskRequest{
			UserId: "test",
		}))
		require.NoError(t, err)
		require.NotEmpty(t, resp.Msg.GetTask())

		resp2, err := client.GenerateTask(context.Background(), connect.NewRequest(&v1.GenerateTaskRequest{
			UserId: "test",
		}))
		require.NoError(t, err)
		require.Equal(t, resp.Msg.GetTask(), resp2.Msg.GetTask())
	})

	t.Run("user_id is empty", func(t *testing.T) {
		_, err := client.GenerateTask(context.Background(), connect.NewRequest(&v1.GenerateTaskRequest{}))
		require.ErrorContains(t, err, "invalid_argument: user_id is required")
		require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}

func TestQuotesHandler_GetQuote(t *testing.T) {
	route, handler := protoconnect.NewLibraryServiceHandler(
		NewQuotesHandler(library.NewLibrary(), pow.NewProofOfWork(1), database.NewTasksStorage()),
	)

	mux := http.NewServeMux()
	mux.Handle(route, handler)

	s := httptest.NewServer(mux)
	t.Cleanup(s.Close)

	client := protoconnect.NewLibraryServiceClient(s.Client(), s.URL)

	t.Run("user_id is empty", func(t *testing.T) {
		_, err := client.GetQuote(context.Background(), connect.NewRequest(&v1.GetQuoteRequest{}))
		require.ErrorContains(t, err, "invalid_argument: user_id is required")
		require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("user_id is empty", func(t *testing.T) {
		_, err := client.GetQuote(context.Background(), connect.NewRequest(&v1.GetQuoteRequest{
			UserId: "test",
		}))
		require.ErrorContains(t, err, "invalid_argument: task is required")
		require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("user_id is empty", func(t *testing.T) {
		_, err := client.GetQuote(context.Background(), connect.NewRequest(&v1.GetQuoteRequest{
			UserId: "test",
			Task:   "test",
		}))
		require.ErrorContains(t, err, "invalid_argument: hash is required")
		require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("task is not generated", func(t *testing.T) {
		_, err := client.GetQuote(context.Background(), connect.NewRequest(&v1.GetQuoteRequest{
			UserId: "test",
			Task:   "test",
			Hash:   "test",
		}))
		require.ErrorContains(t, err, "not_found: task is not generated")
		require.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	})

	t.Run("error - wrong hash", func(t *testing.T) {
		resp, err := client.GenerateTask(context.Background(), connect.NewRequest(&v1.GenerateTaskRequest{
			UserId: "test",
		}))
		require.NoError(t, err)

		_, err = client.GetQuote(context.Background(), connect.NewRequest(&v1.GetQuoteRequest{
			UserId: "test",
			Task:   resp.Msg.GetTask(),
			Hash:   "test",
		}))
		require.ErrorContains(t, err, "unknown: proof of work is invalid")
		require.Equal(t, connect.CodeUnknown, connect.CodeOf(err))
	})

	t.Run("success", func(t *testing.T) {
		resp, err := client.GenerateTask(context.Background(), connect.NewRequest(&v1.GenerateTaskRequest{
			UserId: "test",
		}))
		require.NoError(t, err)

		_, err = client.GetQuote(context.Background(), connect.NewRequest(&v1.GetQuoteRequest{
			UserId: "test",
			Task:   resp.Msg.GetTask(),
			Hash:   "test",
		}))
		require.ErrorContains(t, err, "unknown: proof of work is invalid")
		require.Equal(t, connect.CodeUnknown, connect.CodeOf(err))
	})
}
