package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// ProofOfWork is a service that does proof of work string validation.
type ProofOfWork struct {
	// TODO: difficulty can be a smart one, based on current RPS.
	difficulty int
}

// NewProofOfWork creates a new proof of work service.
func NewProofOfWork(difficulty int) *ProofOfWork {
	return &ProofOfWork{
		difficulty: difficulty,
	}
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
