package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Signature represents a file signature
type Signature struct {
	ID           uuid.UUID // unique identifier for the signature
	FileSize     int64     // size of the file in bytes
	FilePath     string    // path to the file
	LastModified time.Time // last modified timestamp
	CreatedAt    time.Time // timestamp for when the signature was created
	Chunks       []Chunk   // chunks of the file
}

// Print prints the signature to stdout
func (s *Signature) Print() {
	fmt.Println("ID: ", s.ID)
	fmt.Println("File size: ", s.FileSize)
	fmt.Println("File path: ", s.FilePath)
	fmt.Println("Last modified: ", s.LastModified)
	fmt.Println("Created at: ", s.CreatedAt)
	fmt.Println("Number of chunks: ", len(s.Chunks))
	for _, chunk := range s.Chunks {
		fmt.Println("Chunk: ")
		chunk.Print()
	}
}

// ValidateSignature validates the given signature
func (s *Signature) ValidateSignature(other *Signature) bool {
	if s.FileSize != other.FileSize || s.LastModified != other.LastModified || s.CreatedAt != other.CreatedAt || len(s.Chunks) != len(other.Chunks) {
		return false
	}
	for i := range s.Chunks {
		if !s.Chunks[i].ValidateChunk(&other.Chunks[i]) {
			return false
		}
	}
	return true
}
