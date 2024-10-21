package database

import (
	"testing"

	"github.com/kl09/powlibrary/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestTasksStorage(t *testing.T) {
	s := NewTasksStorage()

	s.Add(domain.POWTask{
		Task:       "aaabbb",
		UserID:     "user_id_1",
		Difficulty: 1,
	})

	task := s.GetForUser("user_id_1")
	require.Equal(t, domain.POWTask{
		Task:       "aaabbb",
		UserID:     "user_id_1",
		Difficulty: 1,
	}, *task)

	s.MarkAsUsed(domain.POWTask{
		Task:       "aaabbb",
		UserID:     "user_id_1",
		Difficulty: 1,
	})

	task = s.GetForUser("user_id_1")
	require.Equal(t, domain.POWTask{
		Task:       "aaabbb",
		UserID:     "user_id_1",
		Difficulty: 1,
		IsUsed:     true,
	}, *task)
}
