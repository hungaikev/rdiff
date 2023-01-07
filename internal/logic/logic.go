package logic

import (
	"context"
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

func (l *Logic) Handle(ctx context.Context, file os.File) (*models.Delta, error) {

	return &models.Delta{}, nil

}

func (l *Logic) Shutdown() error {

	return nil
}
