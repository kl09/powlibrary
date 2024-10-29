package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const minDifficulty = 1

// ProofOfWork is a service that does proof of work string validation.
type ProofOfWork struct {
	difficulty int
}

// NewProofOfWork creates a new proof of work service.
func NewProofOfWork(defaultPOWDifficulty int) *ProofOfWork {
	return &ProofOfWork{
		difficulty: defaultPOWDifficulty,
	}
}

// IncreaseDifficulty increases the difficulty of the task.
func (p *ProofOfWork) IncreaseDifficulty() int {
	p.difficulty++

	return p.difficulty
}

// DecreaseDifficulty decreases the difficulty of the task.
func (p *ProofOfWork) DecreaseDifficulty() int {
	if p.difficulty > minDifficulty {
		p.difficulty--
	}

	return p.difficulty
}

// Generate generates a task for POW.
func (p *ProofOfWork) Generate() (string, int, error) {
	rand, err := uuid.NewUUID()
	if err != nil {
		return "", 0, err
	}

	return rand.String(), p.difficulty, nil
}

// Validate validates a proof of work string against a hash.
func (p *ProofOfWork) Validate(task, hash string) error {
	hashToBe := sha256.Sum256([]byte(task + hash))
	result := hex.EncodeToString(hashToBe[:])

	if result[:p.difficulty] != strings.Repeat("0", p.difficulty) {
		return errors.New("proof of work is invalid")
	}

	return nil
}
