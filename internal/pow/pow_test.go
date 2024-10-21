package pow

import (
	"context"
	"testing"

	"github.com/kl09/powlibrary/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestProofOfWork_Validate(t *testing.T) {
	var difficulty = 3

	pow := NewProofOfWork(difficulty)
	task, _, err := pow.Generate()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	got, err := utils.GeneratePOW(ctx, task, difficulty)
	require.NoError(t, err)
	require.NoError(t, pow.Validate(task, got))
}
