// Package apply implements the apply function.
package apply

import (
	"fmt"
	"io"
	"os"

	"github.com/hungaikev/rdiff/internal/shared/models"
)

/*
Apply applies the changes in the Delta struct to the file at the specified path.

1. Creates a new file at the given path using os.Create.
2. Iterates through the Modified field of the Delta struct and writes the modified chunks to the file using file.WriteAt.
3. Iterates through the Added field of the Delta struct and writes the added chunks to the file using file.WriteAt.
4. Returns a nil error value if successful, or returns an error value if there was an error.

*/
// Apply applies the changes described in the given Delta to the original file
func Apply(originalFile string, delta *models.Delta) error {
	// open the original file
	original, err := os.Open(originalFile)
	if err != nil {
		return fmt.Errorf("error opening original file: %w", err)
	}
	defer original.Close()

	// create a new file to store the updated version of the original file
	updated, err := os.Create(originalFile)
	if err != nil {
		return fmt.Errorf("error creating updated file: %w", err)
	}
	defer updated.Close()

	// copy the contents of the original file to the updated file
	if _, err := io.Copy(updated, original); err != nil {
		return fmt.Errorf("error copying original file to updated file: %w", err)
	}

	// apply the changes to the updated file
	for _, chunk := range delta.Added {
		if _, err := updated.Write(chunk.Data); err != nil {
			return fmt.Errorf("error writing added chunk to updated file: %w", err)
		}
	}
	for _, chunk := range delta.Modified {
		if _, err := updated.WriteAt(chunk.Data, chunk.Offset); err != nil {
			return fmt.Errorf("error writing modified chunk to updated file: %w", err)
		}
	}

	return nil
}
