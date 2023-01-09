package logic

import (
	"context"
	"fmt"
	"github.com/hungaikev/rdiff/internal/pkg/apply"
	"github.com/hungaikev/rdiff/internal/pkg/diff"
	"github.com/hungaikev/rdiff/internal/pkg/signature"
	"os"

	"github.com/rs/zerolog"

	"github.com/hungaikev/rdiff/internal/shared/models"
	"github.com/hungaikev/rdiff/internal/store"
)

// Logic defines the business logic for application related operations within this library
type Logic struct {
	log     *zerolog.Logger
	storage store.Storage
}

// New creates an instance of the Logic implementation
func New(log *zerolog.Logger, storage store.Storage) *Logic {
	return &Logic{
		log:     log,
		storage: storage,
	}
}

func (l *Logic) Handle(ctx context.Context, file *os.File) (*models.Delta, error) {
	// step 1: check if file exists in storage
	exists, err := l.storage.FileExists(file.Name())
	if err != nil {
		return nil, fmt.Errorf("error checking if file exists: %w", err)
	}

	var original *models.Signature
	if exists {
		// step 2: retrieve the signature
		original, err = l.storage.GetSignatureForFilename(file.Name())
		if err != nil {
			return nil, fmt.Errorf("error retrieving signature: %w", err)
		}
	}

	// step 3: generate the signature of the new file and store it
	updated, err := signature.GenerateSignature(file)
	if err != nil {
		return nil, fmt.Errorf("error generating signature: %w", err)
	}
	if err := l.storage.Save(updated); err != nil {
		return nil, fmt.Errorf("error saving signature: %w", err)
	}

	var delta *models.Delta
	if exists {
		// step 4: run the Diff method
		delta, err = diff.Diff(original, updated)
		if err != nil {
			return nil, fmt.Errorf("error running Diff: %w", err)
		}
	}

	// step 5: run Apply to synchronize the files
	if err := apply.Apply(original.FilePath, delta); err != nil {
		return nil, fmt.Errorf("error running Apply: %w", err)
	}

	return delta, nil
}

func (l *Logic) Shutdown() error {

	return nil
}
