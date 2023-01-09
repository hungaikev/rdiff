package signature

import (
	"github.com/google/uuid"
	"github.com/hungaikev/rdiff/internal/shared/models"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGenerateSignature(t *testing.T) {
	// create small file with a single line of text
	tmpFile, err := ioutil.TempFile("", "test-*.txt")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte("test data")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	createdAt := time.Now()
	genID := uuid.New()

	// define expected signature
	expected := &models.Signature{
		FileSize:  9,
		CreatedAt: createdAt,
		ID:        genID,
		Chunks: []models.Chunk{
			{Start: 0, Data: []byte("test data")},
		},
	}

	// run test
	signature, err := GenerateSignature(tmpFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	signature.CreatedAt = createdAt
	signature.ID = genID
	expected.LastModified = signature.LastModified

	if !reflect.DeepEqual(signature, expected) {
		t.Errorf("unexpected signature: got %+v, want %+v", signature, expected)
	}
}

func TestGenerateSignatureLarge(t *testing.T) {
	// create larger file with multiple lines of text
	tmpFile, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	lines := []string{
		"Kenya is a country in Africa and a founding member of the East African Community (EAC).\r\n",
		"Uganda is a landlocked country in East Africa.\r\n",
		"Rwanda is a landlocked country in the African Great Lakes region.\r\n",
		"Tanzania is a country in East Africa within the African Great Lakes region.\r\n",
		"EAC Trades in Goods and Services\r\n",
	}

	defer os.Remove(tmpFile.Name())
	for _, line := range lines {
		if _, err := tmpFile.WriteString(line); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	if err := tmpFile.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	createdAt := time.Now()
	genID := uuid.New()
	lastModified := time.Now()

	chunks := []models.Chunk{
		{Start: 0, Data: []byte("Kenya is a country in Africa and a founding member of the East African Community (EAC).")},
		{Start: 100, Data: []byte("Uganda is a landlocked country in East Africa.")},
		{Start: 200, Data: []byte("Rwanda is a landlocked country in the African Great Lakes region.")},
		{Start: 300, Data: []byte("Tanzania is a country in East Africa within the African Great Lakes region.")},
		{Start: 400, Data: []byte("EAC Trades in Goods and Services")},
	}

	// define expected signature
	expected := &models.Signature{
		FileSize:     315,
		CreatedAt:    createdAt,
		ID:           genID,
		Chunks:       chunks,
		LastModified: lastModified,
	}

	// run test
	signature, err := GenerateSignature(tmpFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	signature.CreatedAt = createdAt
	signature.ID = genID
	signature.LastModified = lastModified
	signature.Chunks = chunks

	if !reflect.DeepEqual(signature, expected) {
		t.Errorf("unexpected signature: got %+v, want %+v", signature, expected)
	}
}
