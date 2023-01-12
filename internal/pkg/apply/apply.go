// Package apply implements the apply function.
package apply

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/hungaikev/rdiff/internal/pkg/signature"
	"github.com/hungaikev/rdiff/internal/shared/models"
	"github.com/hungaikev/rdiff/internal/store"
)

// Apply defines the business logic for application related operations within this library
type Apply struct {
	delta   *models.Delta
	storage store.Storage
	log     *zerolog.Logger
	tracer  trace.Tracer
}

// New creates an instance of the Apply implementation
func New(delta *models.Delta, storage store.Storage, log *zerolog.Logger, tracer trace.Tracer) *Apply {
	tracer = otel.Tracer("apply")
	return &Apply{
		delta:   delta,
		storage: storage,
		log:     log,
		tracer:  tracer,
	}
}

// Changes applies the changes described in the given Delta to the original file
func (a *Apply) Changes(ctx context.Context, originalSig *models.Signature) (*models.Signature, error) {
	ctx, span := a.tracer.Start(ctx, "apply.changes")
	defer span.End()

	// open the original file
	original, err := os.Open(originalSig.FilePath)
	if err != nil {
		return &models.Signature{}, fmt.Errorf("error opening original file: %w", err)
	}
	defer original.Close()

	// apply the changes to the original file
	for _, chunk := range a.delta.Added {
		if _, err := original.Write(chunk.Data); err != nil {
			return &models.Signature{}, fmt.Errorf("error writing added chunk to updated file: %w", err)
		}
	}
	for _, chunk := range a.delta.Modified {
		if _, err := original.WriteAt(chunk.Data, chunk.Offset); err != nil {
			return &models.Signature{}, fmt.Errorf("error writing modified chunk to updated file: %w", err)
		}
	}

	a.log.Info().Msgf("Changes applied successfully - original file updated: %s", original.Name())
	a.delta.Print()

	// generate the signature of the updated file.
	newSig, err := signature.Generate(ctx, original)
	if err != nil {
		return &models.Signature{}, fmt.Errorf("error generating signature: %w", err)
	}

	// update the new signature to storage
	saved, err := a.storage.Update(ctx, newSig)
	if err != nil {
		return &models.Signature{}, fmt.Errorf("error saving signature: %w", err)
	}

	a.log.Info().Msgf("Updated signature saved to storage: %s", saved.FilePath)
	saved.Print()

	return saved, nil
}
