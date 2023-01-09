package memory

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/hungaikev/rdiff/internal/shared/models"
)

type Storage struct {
	signatures map[uuid.UUID]*models.Signature
	mu         sync.Mutex
}

func NewStorage() *Storage {
	var mutex sync.Mutex
	return &Storage{
		signatures: make(map[uuid.UUID]*models.Signature),
		mu:         mutex,
	}
}

func (s *Storage) Save(signature *models.Signature) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.signatures[signature.ID] = signature
	return nil
}

func (s *Storage) Get(id uuid.UUID) (*models.Signature, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	signature, ok := s.signatures[id]
	if !ok {
		return nil, fmt.Errorf("signature not found")
	}
	return signature, nil
}

func (s *Storage) ChunkExists(chunk models.Chunk) (bool, error) {
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

func (s *Storage) GetSignatureForChunk(chunk models.Chunk) (*models.Signature, error) {
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

func (s *Storage) FileExists(filename string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sig := range s.signatures {
		if sig.FilePath == filename {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) GetSignatureForFilename(filename string) (*models.Signature, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sig := range s.signatures {
		if sig.FilePath == filename {
			return sig, nil
		}
	}
	return nil, fmt.Errorf("signature not found")
}
