package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Storage interface {
	// Save saves the given signature to storage
	Save(signature *Signature) error

	// Get retrieves the signature with the given ID from storage
	Get(id uuid.UUID) (*Signature, error)

	// ChunkExists checks if the given chunk exists in storage
	ChunkExists(chunk Chunk) (bool, error)

	// GetSignatureForChunk retrieves the signature that contains the given chunk
	GetSignatureForChunk(chunk Chunk) (*Signature, error)
}

type MemoryStorage struct {
	signatures map[uuid.UUID]*Signature
	mu         sync.Mutex
}

func (s *MemoryStorage) Save(signature *Signature) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.signatures[signature.ID] = signature
	return nil
}

func (s *MemoryStorage) Get(id uuid.UUID) (*Signature, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	signature, ok := s.signatures[id]
	if !ok {
		return nil, fmt.Errorf("signature not found")
	}
	return signature, nil
}

func (s *MemoryStorage) ChunkExists(chunk Chunk) (bool, error) {
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

func (s *MemoryStorage) GetSignatureForChunk(chunk Chunk) (*Signature, error) {
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

func (s *MemoryStorage) FileExists(filename string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sig := range s.signatures {
		if sig.FilePath == filename {
			return true, nil
		}
	}
	return false, nil
}
