package database

import (
	"testing"
	"time"

	"github.com/kl09/powlibrary/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestTasksStorage(t *testing.T) {
	now := time.Now()

	s := NewTasksStorage(time.Minute, WithTimeNow(func() time.Time {
		return now
	}))

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
		TTL:        now.Add(time.Minute),
		CreatedAt:  now,
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
		TTL:        now.Add(time.Minute),
		CreatedAt:  now,
		ResolvedAt: &now,
	}, *task)
}

func TestClearCache(t *testing.T) {
	now := time.Now()

	s := NewTasksStorage(500*time.Millisecond, WithTimeNow(func() time.Time {
		return now
	}))

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
		TTL:        now.Add(500 * time.Millisecond),
		CreatedAt:  now,
	}, *task)

	time.Sleep(time.Second)
	s.ClearCache()

	require.Nil(t, s.GetForUser("user_id_1"))
}
