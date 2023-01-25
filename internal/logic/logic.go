// Package logic implements the business logic of the application.
package logic

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	"github.com/hungaikev/rdiff/internal/pkg/apply"
	"github.com/hungaikev/rdiff/internal/pkg/diff"
	"github.com/hungaikev/rdiff/internal/pkg/signature"
	"github.com/hungaikev/rdiff/internal/shared/models"
	"github.com/hungaikev/rdiff/internal/store"
)

// Logic defines the business logic for application related operations within this library
type Logic struct {
	log     *zerolog.Logger
	storage store.Storage
	tracer  trace.Tracer
}

// New creates an instance of the Logic implementation
func New(log *zerolog.Logger, storage store.Storage, tracer trace.Tracer) *Logic {
	return &Logic{
		log:     log,
		storage: storage,
		tracer:  tracer,
	}
}

// Handle handles the application logic for the given library
func (l *Logic) Handle(ctx context.Context, file *os.File) (*models.Delta, error) {
	ctx, span := l.tracer.Start(ctx, "logic.Handle")
	defer span.End()

	l.log.Info().Msgf("handling diff logic: %s", file.Name())

	// step 1: check if file exists in storage
	exists, err := l.storage.FileExists(ctx, file.Name())
	if err != nil {
		return nil, fmt.Errorf("error checking if file exists: %w", err)
	}

	// step 2: Signature generation
	var original *models.Signature
	var delta *models.Delta
	if exists {
		// step 2: retrieve the signature if file exists
		original, err = l.storage.GetSignatureForFilename(ctx, file.Name())
		if err != nil {
			return nil, fmt.Errorf("error retrieving signature: %w", err)
		}

		// step 2.1 generate the signature of the updated file.
		updated, err := signature.Generate(ctx, file, l.log)
		if err != nil {
			return nil, fmt.Errorf("error generating signature: %w", err)
		}

		// step 2.2: run the Diff method
		delta, err = diff.Compare(ctx, original, updated)
		if err != nil {
			return nil, fmt.Errorf("error running Diff: %w", err)
		}

		// instantiate a new apply
		apply := apply.New(delta, l.storage, l.log, l.tracer)

		// step 2.3: run apply to synchronize the files
		newSig, err := apply.Changes(ctx, original)
		if err != nil {
			return nil, fmt.Errorf("error applying changes: %w", err)
		}

		original = newSig

		return delta, nil

	}

	// step 3: generate the signature of the new file and store it
	updated, err := signature.Generate(ctx, file, l.log)
	if err != nil {
		return nil, fmt.Errorf("error generating signature: %w", err)
	}
	l.log.Info().Msgf("generated signature for file: %s", file.Name())

	saved, err := l.storage.Save(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("error saving signature: %w", err)
	}

	original = saved

	return delta, nil
}
