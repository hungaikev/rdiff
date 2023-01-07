package main

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	// define test case
	original := &Signature{
		FileSize: 100,
		Chunks: []Chunk{
			{Start: 0, Data: []byte("chunk 1")},
			{Start: 8, Data: []byte("chunk 2")},
			{Start: 16, Data: []byte("chunk 3")},
		},
	}
	updated := &Signature{
		FileSize: 100,
		Chunks: []Chunk{
			{Start: 0, Data: []byte("chunk 1")},
			{Start: 8, Data: []byte("chunk 2 modified")},
			{Start: 16, Data: []byte("chunk 3")},
		},
	}
	expected := &Delta{
		Added: []Chunk{},
		Modified: []Chunk{
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

	added := []Chunk{
		{Start: 24, Data: []byte("chunk 4")},
		{Start: 32, Data: []byte("chunk 5")},
	}

	modified := []Chunk{
		{Start: 8, Data: []byte("chunk 2 modified")},
	}

	original := &Signature{
		FileSize: 100,
		Chunks: []Chunk{
			{Start: 0, Data: []byte("chunk 1")},
			{Start: 8, Data: []byte("chunk 2")},
			{Start: 16, Data: []byte("chunk 3")},
		},
	}
	updated := &Signature{
		FileSize: 150,
		Chunks: []Chunk{
			{Start: 0, Data: []byte("chunk 1")},
			{Start: 8, Data: []byte("chunk 2 modified")},
			{Start: 16, Data: []byte("chunk 3")},
			{Start: 24, Data: []byte("chunk 4")},
			{Start: 32, Data: []byte("chunk 5")},
		},
	}
	expected := &Delta{
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
