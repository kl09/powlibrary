package domain

import (
	"context"
	"time"
)

// POWTask is a proof of work task which was given to a user.
type POWTask struct {
	Task       string
	Difficulty int
	UserID     string

	CreatedAt  time.Time
	ResolvedAt *time.Time
	TTL        time.Time
}

// TasksStorage is a service that stores tasks that were given for a user.
// TODO: we can use PG here.
type TasksStorage interface {
	// GetForUser returns a task for a user.
	GetForUser(userID string) *POWTask
	// Add adds a task for a user.
	Add(task POWTask)
	// MarkAsUsed marks a task as used.
	MarkAsUsed(task POWTask)
}

// LibraryService is a service that provides quotes.
type LibraryService interface {
	// GetRandomQuote returns a random quote.
	GetRandomQuote(ctx context.Context) string
}

// POWService is a service that does proof of work string validation.
type POWService interface {
	Generate() (string, int, error)
	// Validate validates a proof of work string against a hash.
	Validate(code, hash string) error

	// IncreaseDifficulty increases the difficulty of the task.
	IncreaseDifficulty() int
	// DecreaseDifficulty decreases the difficulty of the task.
	DecreaseDifficulty() int
}
