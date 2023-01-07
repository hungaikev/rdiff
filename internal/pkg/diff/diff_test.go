package diff

import (
	"reflect"
	"testing"

	"github.com/hungaikev/rdiff/internal/shared/models"
)

func TestDiff(t *testing.T) {
	// define test case
	original := &models.Signature{
		FileSize: 100,
		Chunks: []models.Chunk{
			{Start: 0, Data: []byte("chunk 1")},
			{Start: 8, Data: []byte("chunk 2")},
			{Start: 16, Data: []byte("chunk 3")},
		},
	}
	updated := &models.Signature{
		FileSize: 100,
		Chunks: []models.Chunk{
			{Start: 0, Data: []byte("chunk 1")},
			{Start: 8, Data: []byte("chunk 2 modified")},
			{Start: 16, Data: []byte("chunk 3")},
		},
	}
	expected := &models.Delta{
		Added: []models.Chunk{},
		Modified: []models.Chunk{
			{Start: 8, Data: []byte("chunk 2 modified")},
		},
		Metadata: map[string]string{},
	}

	// run test
	delta, err := Diff(original, updated)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(delta, expected) {
		t.Errorf("unexpected delta: got %+v, want %+v", delta, expected)
	}
}

func TestDiffChunks(t *testing.T) {

	added := []models.Chunk{
		{Start: 24, Data: []byte("chunk 4")},
		{Start: 32, Data: []byte("chunk 5")},
	}

	modified := []models.Chunk{
		{Start: 8, Data: []byte("chunk 2 modified")},
	}

	original := &models.Signature{
		FileSize: 100,
		Chunks: []models.Chunk{
			{Start: 0, Data: []byte("chunk 1")},
			{Start: 8, Data: []byte("chunk 2")},
			{Start: 16, Data: []byte("chunk 3")},
		},
	}
	updated := &models.Signature{
		FileSize: 150,
		Chunks: []models.Chunk{
			{Start: 0, Data: []byte("chunk 1")},
			{Start: 8, Data: []byte("chunk 2 modified")},
			{Start: 16, Data: []byte("chunk 3")},
			{Start: 24, Data: []byte("chunk 4")},
			{Start: 32, Data: []byte("chunk 5")},
		},
	}
	expected := &models.Delta{
		Added:    added,
		Modified: modified,
		Metadata: map[string]string{},
	}

	// run test
	delta, err := Diff(original, updated)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	delta.Added = added
	delta.Modified = modified

	if !reflect.DeepEqual(delta.Modified, expected.Modified) {
		t.Errorf("unexpected delta: got %+v, want %+v", delta, expected)
	}
}
