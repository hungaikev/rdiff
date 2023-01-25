// Package signature implements functions for generating and comparing file signatures.
package signature

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/hungaikev/rdiff/internal/pkg/chunks"
	"github.com/hungaikev/rdiff/internal/shared/models"
)

var tracer = otel.Tracer("signature")

// Generate generates a new signature for the given file and returns it
func Generate(ctx context.Context, file *os.File, log *zerolog.Logger) (*models.Signature, error) {
	ctx, span := tracer.Start(ctx, "signature.Generate")
	defer span.End()

	// get file information
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file information: %w", err)
	}

	// generate chunks
	chunks, err := chunks.Generate(ctx, file, log)
	if err != nil {
		return nil, fmt.Errorf("unable to generate chunks: %w", err)
	}

	// create a new signature
	signature := &models.Signature{
		ID:           uuid.New(),
		FileSize:     fileInfo.Size(),
		FilePath:     file.Name(),
		LastModified: fileInfo.ModTime(),
		CreatedAt:    time.Now().UTC(),
		Chunks:       chunks,
	}

	log.Info().Msgf("generated signature for file %s", file.Name())
	signature.Print()

	return signature, nil
}
