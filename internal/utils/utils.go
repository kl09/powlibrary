package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GeneratePOW(ctx context.Context, task string, difficulty int) (string, error) {
	var (
		result string
		found  bool
		randID uuid.UUID
		err    error
	)
	for {
		select {
		case <-ctx.Done():
			return "", errors.New("context done")
		default:

		}

		randID, err = uuid.NewUUID()
		if err != nil {
			return "", fmt.Errorf("generate uuid: %w", err)
		}

		hash := sha256.Sum256([]byte(task + randID.String()))
		result = hex.EncodeToString(hash[:])
		if result[:difficulty] == strings.Repeat("0", difficulty) {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("can't find valid hash")
	}

	return randID.String(), nil
}
