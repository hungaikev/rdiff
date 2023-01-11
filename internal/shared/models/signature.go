package models

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/hungaikev/rdiff/internal/pkg/rolling"
)

const chunkSize = 8192 // size of each chunk in bytes

// Signature represents a file signature
type Signature struct {
	ID           uuid.UUID // unique identifier for the signature
	FileSize     int64     // size of the file in bytes
	FilePath     string    // path to the file
	LastModified time.Time // last modified timestamp
	CreatedAt    time.Time // timestamp for when the signature was created
	Chunks       []Chunk   // chunks of the file
}

// NewGeneratedSignature generates a new signature for the given file and returns it
func NewGeneratedSignature(file *os.File) (*Signature, error) {
	// get file information
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file information: %w", err)
	}

	// generate chunks
	chunks, err := generateChunks(file)
	if err != nil {
		return nil, fmt.Errorf("unable to generate chunks: %w", err)
	}

	// create a new signature
	signature := &Signature{
		ID:           uuid.New(),
		FileSize:     fileInfo.Size(),
		FilePath:     file.Name(),
		LastModified: fileInfo.ModTime(),
		CreatedAt:    time.Now().UTC(),
		Chunks:       chunks,
	}

	return signature, nil
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

// generateChunks reads the given file chunk by chunk and returns a slice of Chunk structs
func generateChunks(file *os.File) ([]Chunk, error) {
	// create a buffer to read the file chunk by chunk
	buf := make([]byte, chunkSize)

	// create a slice to store the chunks
	chunks := make([]Chunk, 0)

	// initialize the rolling hash value to 0
	rollingHash := uint64(0)

	// initialize the offset to 0
	offset := int64(0)

	// read the file chunk by chunk
	for {
		// read a chunk of data
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			// if there was an error other than EOF, return it
			return nil, err
		}
		if n == 0 {
			// if no data was read, we've reached the end of the file
			break
		}

		// calculate the rolling hash value for the chunk
		rollingHash = rolling.Hash(buf[:n])

		// update the offset
		offset += int64(n)

		// create a new chunk with the data and rolling hash value
		chunk := Chunk{
			Data:   buf[:n],
			Hash:   rollingHash,
			Offset: offset,
			Length: int64(n),
		}

		// add the chunk to the slice
		chunks = append(chunks, chunk)

	}

	return chunks, nil
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
