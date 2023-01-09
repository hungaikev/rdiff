package tests

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hungaikev/rdiff/internal/pkg/apply"
	"github.com/hungaikev/rdiff/internal/pkg/diff"
	"github.com/hungaikev/rdiff/internal/pkg/fileio"
	"github.com/hungaikev/rdiff/internal/pkg/signature"
)

func TestEndToEnd(t *testing.T) {

	tmpOriginal, err := fileio.OpenFile("testdata/tmp-original.txt")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	defer tmpOriginal.Close()

	// generate signature for original file
	original, err := signature.GenerateSignature(tmpOriginal)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tmpUpdated, err := fileio.OpenFile("testdata/tmp-updated.txt")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	defer tmpUpdated.Close()

	// generate signature for updated file
	updated, err := signature.GenerateSignature(tmpUpdated)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

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

	assert.Contains(t, string(b), "updated")
}
