package models

import (
	"fmt"
)

// Delta represents a delta between two files
type Delta struct {
	Added    []Chunk           // list of chunks that have been added
	Modified []Chunk           // list of chunks that have been modified
	Metadata map[string]string // metadata needed to apply the changes to the original signature
}

// Print prints the delta to stdout
func (d *Delta) Print() {
	fmt.Println("Added chunks:")
	for _, chunk := range d.Added {
		fmt.Printf("  start: %d, data: %s\n", chunk.Start, chunk.Data)
	}
	fmt.Println("Modified chunks:")
	for _, chunk := range d.Modified {
		fmt.Printf("  start: %d, data: %s\n", chunk.Start, chunk.Data)
	}
	fmt.Println("Metadata:")
	for key, value := range d.Metadata {
		fmt.Printf("  %s: %s\n", key, value)
	}
}
