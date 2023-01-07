package store

import (
	"github.com/google/uuid"

	"github.com/hungaikev/rdiff/internal/shared/models"
)

type Storage interface {
	// Save saves the given signature to storage
	Save(signature *models.Signature) error

	// Get retrieves the signature with the given ID from storage
	Get(id uuid.UUID) (*models.Signature, error)

	// ChunkExists checks if the given chunk exists in storage
	ChunkExists(chunk models.Chunk) (bool, error)

	// GetSignatureForChunk retrieves the signature that contains the given chunk
	GetSignatureForChunk(chunk models.Chunk) (*models.Signature, error)
}
