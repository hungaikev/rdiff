package main

import "fmt"

// Chunk represents a chunk of data
type Chunk struct {
	Start int64  // starting position of the chunk in the file
	Data  []byte // actual data for the chunk
}

// Print prints the chunk to stdout
func (c *Chunk) Print() {
	fmt.Printf("start: %d, data: %s\n", c.Start, c.Data)
}
