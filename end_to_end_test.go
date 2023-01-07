package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestEndToEnd(t *testing.T) {
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

	// generate signature for original file
	original, err := GenerateSignature(tmpOriginal.Name())
	if err != nil {
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

	// generate signature for updated file
	updated, err := GenerateSignature(tmpUpdated.Name())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// compute delta
	delta, err := Diff(original, updated)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// apply delta to original file
	if err := Apply(tmpOriginal.Name(), delta); err != nil {
		t.Errorf("un	expected error: %v", err)
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
