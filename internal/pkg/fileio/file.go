// Package fileio implements file I/O functions.
package fileio

import "os"

// OpenFile opens a file and returns a pointer to it, along with an error value.
// If there is an error opening the file, the error value will be non-nil.
// Otherwise, the error value will be nil and the pointer to the file can be used
// to read or write to the file.
func OpenFile(filename string) (*os.File, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
