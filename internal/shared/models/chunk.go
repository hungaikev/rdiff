package models

import (
	"bytes"
	"fmt"
)

// Chunk represents a chunk of data
type Chunk struct {
	Data   []byte // the data of the chunk
	Hash   uint64 // the rolling hash value of the data
	Offset int64  // the starting position of the chunk in the file
	Length int64  // the length of the chunk
}

// ValidateChunk validates chunks
func (c *Chunk) ValidateChunk(other *Chunk) bool {
	if c.Length != other.Length {
		return false
	}
	if c.Offset != other.Offset {
		return false
	}
	if c.Hash != other.Hash {
		return false
	}
	if !bytes.Equal(c.Data, other.Data) {
		return false
	}
	return true
}

// Print prints the chunk
func (c *Chunk) Print() {
	fmt.Printf("Data: %v\n", c.Data)
	fmt.Printf("Hash: %d\n", c.Hash)
	fmt.Printf("Offset: %d\n", c.Offset)
	fmt.Printf("Length: %d\n", c.Length)
}
