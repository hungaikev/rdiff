package apply

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hungaikev/rdiff/internal/shared/models"
)

func TestApply(t *testing.T) {
	// create temporary original file
	tmpOriginal, err := os.Create("tmp-original.txt")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer os.Remove(tmpOriginal.Name())
	if _, err := tmpOriginal.Write([]byte("original data")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := tmpOriginal.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// create temporary updated file
	tmpUpdated, err := os.Create("tmp-updated.txt")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer os.Remove(tmpUpdated.Name())
	if _, err := tmpUpdated.Write([]byte("updated data")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := tmpUpdated.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// create delta
	delta := &models.Delta{
		Added: []models.Chunk{
			{Start: 0, Data: []byte("updated ")},
			{Start: 8, Data: []byte("data")},
		},
		Modified: []models.Chunk{
			{Start: 0, Data: []byte("updated ")},
			{Start: 8, Data: []byte("data")},
		},
		Metadata: map[string]string{"key": "value"},
	}

	// run test
	if err := Apply(tmpOriginal.Name(), delta); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// check if original file has been updated
	b, err := ioutil.ReadFile(tmpOriginal.Name())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(b) != "updated data" {
		t.Errorf("unexpected file content: got %q, want %q", string(b), "updated data")
	}
}
