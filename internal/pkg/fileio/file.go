// Package fileio implements file I/O functions.
package fileio

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("fileio")

// OpenFile opens a file and returns a pointer to it, along with an error value.
// If there is an error opening the file, the error value will be non-nil.
// Otherwise, the error value will be nil and the pointer to the file can be used
// to read or write to the file.
func OpenFile(ctx context.Context, filename string) (*os.File, error) {
	ctx, span := tracer.Start(ctx, "fileio.OpenFile")
	defer span.End()

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// WriteToFile writes to a file
func WriteToFile(ctx context.Context, filePath string, data string) error {
	ctx, span := tracer.Start(ctx, "fileio.WriteToFile")
	defer span.End()

	// Open the file in write mode
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	// Write the data to the file
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
