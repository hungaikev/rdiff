// Package diff implements the diff function.
package diff

import (
	"bytes"
	"context"

	"go.opentelemetry.io/otel"

	"github.com/hungaikev/rdiff/internal/shared/models"
)

var tracer = otel.Tracer("diff")

/*
Compare compares the original and updated signatures and returns a Delta struct containing the differences between the two.

1. Creates a new Delta struct and initializes its fields with empty slices and a map.
2. Compares the number of chunks in the original and updated signatures. If they are not the same, it assumes that the entire file has been modified and sets the Modified field of the Delta struct to the Chunks field of the updated signature.
3. If the number of chunks is the same, it compares the rolling hash values of each chunk using a for loop. If the hash values are different, it assumes that the chunk has been modified and adds it to the Modified field of the Delta struct.
4. Returns the Delta pointer and a nil error value if successful, or returns a nil pointer and an error value if there was an error.
*/
func Compare(ctx context.Context, original, updated *models.Signature) (*models.Delta, error) {
	ctx, span := tracer.Start(ctx, "diff.Compare")
	defer span.End()

	// create a new delta
	delta := &models.Delta{
		Added:    make([]models.Chunk, 0),
		Modified: make([]models.Chunk, 0),
		Metadata: make(map[string]string),
	}

	// check if the original and updated signatures have the same number of chunks
	if len(original.Chunks) != len(updated.Chunks) {
		// if not, the entire file has been modified
		delta.Modified = updated.Chunks
		return delta, nil
	}

	// compare the rolling hash values of each chunk
	for i := range original.Chunks {

		result := bytes.Compare(original.Chunks[i].Data, updated.Chunks[i].Data)

		if result != 0 {
			// if the hash values are different, the chunk has been modified
			delta.Modified = append(delta.Modified, updated.Chunks[i])
		}
	}

	return delta, nil
}
