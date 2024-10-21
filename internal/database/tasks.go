package database

import (
	"sync"

	"github.com/kl09/powlibrary/internal/domain"
)

// TasksStorage is a service that stores tasks that were given for a user.
// TODO: we can change this memory based implementation to PG.
type TasksStorage struct {
	// storage for task: map[user_id]domain.POWTask
	storage map[string]domain.POWTask

	mw sync.RWMutex
}

// NewTasksStorage creates a new task storage service.
func NewTasksStorage() *TasksStorage {
	return &TasksStorage{
		storage: make(map[string]domain.POWTask),
	}
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

	s.storage[task.UserID] = task
}

func (s *TasksStorage) MarkAsUsed(task domain.POWTask) {
	s.mw.Lock()
	defer s.mw.Unlock()

	t, ok := s.storage[task.UserID]
	if ok {
		s.storage[task.UserID] = domain.POWTask{
			Task:       t.Task,
			UserID:     t.UserID,
			Difficulty: task.Difficulty,
			IsUsed:     true,
		}
	}
}
