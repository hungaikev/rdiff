// Package memory implements the storage interface for the memory storage
package memory

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/hungaikev/rdiff/internal/shared/models"
)

// Storage is the memory storage implementation
type Storage struct {
	signatures map[uuid.UUID]*models.Signature
	mu         sync.Mutex
	log        *zerolog.Logger
	tracer     trace.Tracer
}

// New creates a new memory storage
func New(log *zerolog.Logger, tracer trace.Tracer) *Storage {
	var mutex sync.Mutex
	tracer = otel.Tracer("memory")

	return &Storage{
		signatures: make(map[uuid.UUID]*models.Signature),
		mu:         mutex,
		log:        log,
		tracer:     tracer,
	}
}

// Save saves the file signature
func (s *Storage) Save(ctx context.Context, signature *models.Signature) (*models.Signature, error) {
	ctx, span := s.tracer.Start(ctx, "memory.Save")
	defer span.End()

	s.mu.Lock()
	defer s.mu.Unlock()

	signature.ID = uuid.New()
	signature.CreatedAt = time.Now()

	s.signatures[signature.ID] = signature

	return signature, nil
}

// Get returns the signature for the given id
func (s *Storage) Get(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	ctx, span := s.tracer.Start(ctx, "memory.Get")
	defer span.End()

	s.mu.Lock()
	defer s.mu.Unlock()

	signature, ok := s.signatures[id]
	if !ok {
		return nil, fmt.Errorf("signature not found")
	}
	return signature, nil
}

// Update updates the given signature
func (s *Storage) Update(ctx context.Context, signature *models.Signature) (*models.Signature, error) {
	ctx, span := s.tracer.Start(ctx, "memory.Update")
	defer span.End()

	s.mu.Lock()
	defer s.mu.Unlock()

	signature.LastModified = time.Now()

	s.signatures[signature.ID] = signature

	return signature, nil
}

// ChunkExists checks if a chunk exists
func (s *Storage) ChunkExists(ctx context.Context, chunk models.Chunk) (bool, error) {
	ctx, span := s.tracer.Start(context.Background(), "memory.ChunkExists")
	defer span.End()

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sig := range s.signatures {
		for _, c := range sig.Chunks {
			if bytes.Equal(c.Data, chunk.Data) {
				return true, nil
			}
		}
	}
	return false, nil
}

// GetSignatureForChunk returns the signature that contains the given chunk
func (s *Storage) GetSignatureForChunk(ctx context.Context, chunk models.Chunk) (*models.Signature, error) {
	ctx, span := s.tracer.Start(ctx, "memory.GetSignatureForChunk")
	defer span.End()

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sig := range s.signatures {
		for _, c := range sig.Chunks {
			if bytes.Equal(c.Data, chunk.Data) {
				return sig, nil
			}
		}
	}
	return nil, fmt.Errorf("signature not found")
}

// FileExists checks if a file exists
func (s *Storage) FileExists(ctx context.Context, filename string) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "memory.FileExists")
	defer span.End()

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sig := range s.signatures {
		if sig.FilePath == filename {
			return true, nil
		}
	}
	return false, nil
}

// GetSignatureForFilename returns the signature for the given file
func (s *Storage) GetSignatureForFilename(ctx context.Context, filename string) (*models.Signature, error) {
	ctx, span := s.tracer.Start(ctx, "memory.GetSignatureForFilename")
	defer span.End()

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sig := range s.signatures {
		if sig.FilePath == filename {
			return sig, nil
		}
	}
	return nil, fmt.Errorf("signature not found")
}
