// Package store is a package that contains the storage interface used by the application.
package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/hungaikev/rdiff/internal/shared/models"
)

type Storage interface {
	// Save saves the given signature to storage
	Save(ctx context.Context, signature *models.Signature) (*models.Signature, error)

	// Get retrieves the signature with the given ID from storage
	Get(ctx context.Context, id uuid.UUID) (*models.Signature, error)

	// Update updates the signature with the given ID in storage
	Update(ctx context.Context, signature *models.Signature) (*models.Signature, error)

	// ChunkExists checks if the given chunk exists in storage
	ChunkExists(ctx context.Context, chunk models.Chunk) (bool, error)

	// GetSignatureForChunk retrieves the signature that contains the given chunk
	GetSignatureForChunk(ctx context.Context, chunk models.Chunk) (*models.Signature, error)

	// FileExists checks if the given file exists in storage
	FileExists(ctx context.Context, filename string) (bool, error)

	// GetSignatureForFilename retrieves the signature that contains the given file
	GetSignatureForFilename(ctx context.Context, filename string) (*models.Signature, error)
}
