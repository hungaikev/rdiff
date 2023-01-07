package main

import (
	"os"
)

/*
Apply applies the changes in the Delta struct to the file at the specified path.

1. Creates a new file at the given path using os.Create.
2. Iterates through the Modified field of the Delta struct and writes the modified chunks to the file using file.WriteAt.
3. Iterates through the Added field of the Delta struct and writes the added chunks to the file using file.WriteAt.
4. Returns a nil error value if successful, or returns an error value if there was an error.

*/

func Apply(path string, delta *Delta) error {
	// create a new file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the modified chunks
	for _, c := range delta.Modified {
		if _, err := file.WriteAt(c.Data, c.Start); err != nil {
			return err
		}
	}

	// write the added chunks
	for _, c := range delta.Added {
		if _, err := file.WriteAt(c.Data, c.Start); err != nil {
			return err
		}
	}

	return nil
}
