package tests

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hungaikev/rdiff/internal/pkg/apply"
	"github.com/hungaikev/rdiff/internal/pkg/diff"
	"github.com/hungaikev/rdiff/internal/pkg/signature"
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
	defer tmpOriginal.Close()

	tmpOriginal.Stat()

	// generate signature for original file
	original, err := signature.GenerateSignature(tmpOriginal)
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
	defer tmpUpdated.Close()

	// generate signature for updated file
	updated, err := signature.GenerateSignature(tmpUpdated)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	original.Print()

	// compute delta
	delta, err := diff.Diff(original, updated)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// apply delta to original file
	if err := apply.Apply(tmpOriginal.Name(), delta); err != nil {
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
