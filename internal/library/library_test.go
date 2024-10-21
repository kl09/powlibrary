package library

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLibrary_GetRandomQuote(t *testing.T) {
	s := NewLibrary()

	q1 := s.GetRandomQuote(context.Background())
	require.NotEmpty(t, q1)

	q2 := s.GetRandomQuote(context.Background())
	require.NotEmpty(t, q2)

	q3 := s.GetRandomQuote(context.Background())
	require.NotEmpty(t, q3)

	require.NotEqualValues(t, q1, []string{q2, q3})
}
