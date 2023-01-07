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
	fmt.Println("File size:", s.FileSize)
	fmt.Println("File path:", s.FilePath)
	fmt.Println("Created at:", s.CreatedAt)
	fmt.Println("ID:", s.ID)
	fmt.Println("Chunks:")
	for _, chunk := range s.Chunks {
		fmt.Printf("  start: %d, data: %s\n", chunk.Start, chunk.Data)
	}
}
