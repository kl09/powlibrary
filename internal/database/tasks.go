package database

import (
	"sync"
	"time"

	"github.com/kl09/powlibrary/internal/domain"
)

// TasksStorage is a service that stores tasks that were given for a user.
// TODO: we can change this memory based implementation to PG.
type TasksStorage struct {
	// storage for task: map[user_id]domain.POWTask
	storage       map[string]domain.POWTask
	cacheDuration time.Duration
	nowFn         func() time.Time

	mw sync.RWMutex
}

// NewTasksStorage creates a new task storage service.
func NewTasksStorage(cacheDuration time.Duration, opts ...Option) *TasksStorage {
	s := &TasksStorage{
		storage:       make(map[string]domain.POWTask),
		cacheDuration: cacheDuration,
		nowFn:         time.Now,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Option func(c *TasksStorage)

func WithTimeNow(nowFn func() time.Time) Option {
	return func(s *TasksStorage) {
		s.nowFn = nowFn
	}
}

// ClearCache clears expired tasks.
func (s *TasksStorage) ClearCache() {
	s.mw.Lock()
	defer s.mw.Unlock()

	for k, v := range s.storage {
		if v.TTL.Before(time.Now()) {
			delete(s.storage, k)
		}
	}
}

func (s *TasksStorage) AvgTimeToResolve() float64 {
	s.mw.Lock()
	defer s.mw.Unlock()

	avgSeconds := 0.0
	i := 0
	for k, v := range s.storage {
		if v.ResolvedAt != nil {
			avgSeconds = avgSeconds + v.ResolvedAt.Sub(v.CreatedAt).Seconds()
			i++
			delete(s.storage, k)
		}
	}

	return avgSeconds / float64(i)
}

// GetForUser returns a task for a user.
func (s *TasksStorage) GetForUser(userID string) *domain.POWTask {
	s.mw.RLock()
	defer s.mw.RUnlock()

	t, ok := s.storage[userID]
	if ok {
		return &t
	}

	return nil
}

// Add adds a task for a user.
func (s *TasksStorage) Add(task domain.POWTask) {
	s.mw.Lock()
	defer s.mw.Unlock()

	task.TTL = s.nowFn().Add(s.cacheDuration)
	task.CreatedAt = s.nowFn()
	s.storage[task.UserID] = task
}

func (s *TasksStorage) MarkAsUsed(task domain.POWTask) {
	s.mw.Lock()
	defer s.mw.Unlock()

	t, ok := s.storage[task.UserID]
	if ok {
		now := s.nowFn()
		t.ResolvedAt = &now
		s.storage[t.UserID] = t
	}
}
