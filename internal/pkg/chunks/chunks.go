package chunks

import (
	"context"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"io"
	"os"

	"github.com/hungaikev/rdiff/internal/pkg/rolling"
	"github.com/hungaikev/rdiff/internal/shared/models"
)

var tracer = otel.Tracer("chunks")

const chunkSize = 8192 // size of each chunk in bytes

// Generate reads the given file chunk by chunk and returns a slice of Chunk structs
func Generate(ctx context.Context, file *os.File, log *zerolog.Logger) ([]models.Chunk, error) {
	ctx, span := tracer.Start(ctx, "chunks.Generate")
	defer span.End()

	// create a buffer to read the file chunk by chunk
	buf := make([]byte, chunkSize)

	// create a slice to store the chunks
	chunks := make([]models.Chunk, 0)

	// initialize the rolling hash value to 0
	rollingHash := uint64(0)

	// initialize the offset to 0
	offset := int64(0)

	// read the file chunk by chunk
	for {
		// read a chunk of data
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			// if there was an error other than EOF, return it
			return nil, err
		}
		if n == 0 {
			// if no data was read, we've reached the end of the file
			break
		}

		// calculate the rolling hash value for the chunk
		rollingHash = rolling.Hash(ctx, buf[:n])

		// update the offset
		offset += int64(n)

		// create a new chunk with the data and rolling hash value
		chunk := models.Chunk{
			Data:   buf[:n],
			Hash:   rollingHash,
			Offset: offset,
			Length: int64(n),
		}

		// add the chunk to the slice
		chunks = append(chunks, chunk)

	}

	log.Info().Msgf("generated %d chunks for file %s", len(chunks), file.Name())

	return chunks, nil
}
