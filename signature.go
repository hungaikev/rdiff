package main

import (
	"crypto/sha1"
	"fmt"
	"io"
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

const chunkSize = 8192 // size of each chunk in bytes

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

/*
GenerateSignature generates a signature for the file at the given path.

1. Opens the file at the given path using os.Open.
2. Retrieves the file information using file.Stat.
3. Creates a new Signature struct and initializes its fields with the file size, last modified timestamp, and current timestamp.
4. Creates a new SHA1 hash value using sha1.New.
5. Reads the file chunk by chunk using a for loop.
6. For each chunk, it calculates the rolling hash value by adding the chunk data to the hash value using hash.Write, creates a new Chunk struct, and adds it to the Chunks slice of the Signature struct.
7. Returns the Signature pointer and a nil error value if successful, or returns a nil pointer and an error value if there was an error.
*/
func GenerateSignature(path string) (*Signature, error) {
	// open the file
	file, err := OpenFile(path)
	defer file.Close()

	// get file info
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// create a new signature
	sig := &Signature{
		ID:           uuid.New(),
		FileSize:     info.Size(),
		LastModified: info.ModTime(),
		CreatedAt:    time.Now(),
		Chunks:       make([]Chunk, 0),
	}

	// create a rolling hash value
	hash := sha1.New()

	// read the file chunk by chunk
	buf := make([]byte, chunkSize)
	for {
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		// add the chunk data to the rolling hash
		hash.Write(buf[:n])

		// create a new chunk
		c := Chunk{
			Start: int64(len(sig.Chunks)) * chunkSize,
			Data:  buf[:n],
		}

		// add the chunk to the signature
		sig.Chunks = append(sig.Chunks, c)
	}

	return sig, nil
}
