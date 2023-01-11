// Package fileio implements file I/O functions.
package fileio

import (
	"bufio"
	"os"
	"strings"
)

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

// ReadFile reads a file and returns the contents as a string, along with an error value.
func ReadFile(file *os.File) (string, error) {
	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Use the Scan function to iterate through the lines of the file
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Return the contents of the file as a single string
	return strings.Join(lines, "\n"), scanner.Err()
}
