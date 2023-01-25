package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"math/rand"
	"os"
	"time"

	"github.com/hungaikev/rdiff/internal/logic"
	"github.com/hungaikev/rdiff/internal/pkg/fileio"
	"github.com/hungaikev/rdiff/internal/store/memory"
)

const (
	// LogStrKeyModule is for use with the logger as a key to specify the module name.
	LogStrKeyModule = "module"
	// LogStrKeyService is for use with the logger as a key to specify the service name.
	LogStrKeyService = "service"

	filePath = "cmd/filestorage/files/test.txt"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	z := zerolog.New(os.Stderr).With().Str(LogStrKeyService, "rdiff").Timestamp().Logger().With().Caller().Logger()
	mainLog := z.With().Str(LogStrKeyModule, "main").Logger()
	mainLog.Info().Msg("starting program...")

	if err := run(&mainLog); err != nil {
		mainLog.Info().Msgf("main: error %s:", err.Error())
		os.Exit(1)
	}
}

func run(log *zerolog.Logger) error {
	log.Info().Msg("Welcome to the rdiff program :)")

	// =========================================================================
	// Configuration

	var cfg struct {
		conf.Version
	}
	cfg.Version.Build = build
	cfg.Version.Desc = "Hungai' Interview Solution"

	// =========================================================================
	// App Starting

	expvar.NewString("build").Set(build)
	log.Info().Msgf("Started: Application initializing: version %q", build)
	defer log.Info().Msg("Completed")

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			log.Info().Msg(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Info().Msgf("Config:\n%v\n", out)

	// =========================================================================
	// Tracer

	var tracer = otel.Tracer("main")

	// =========================================================================
	// Storage

	storage := memory.New(log, tracer)

	// =========================================================================
	// Logic

	rdiffLogic := logic.New(log, storage, tracer)

	// =========================================================================
	// Open File

	ctx := context.Background()

	log.Info().Msg("Opening file...")

	file, err := fileio.OpenFile(ctx, filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()

	// =========================================================================
	// Handle File
	_, err = rdiffLogic.Handle(ctx, file)
	if err != nil {
		return fmt.Errorf("error handling file: %w", err)
	}

	// =========================================================================

	log.Info().Msgf("File %s has been processed successfully", filePath)

	// =========================================================================

	log.Info().Msgf("Writing to file %s", filePath)

	if err := fileio.WriteToFile(ctx, filePath, generateRandomString(100)); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	file2, err := fileio.OpenFile(ctx, filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer file2.Close()

	_, err = rdiffLogic.Handle(ctx, file2)
	if err != nil {
		return fmt.Errorf("error handling file 2: %w", err)
	}

	return nil
}

func generateRandomString(length int) string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// Define a set of characters to use in the random string
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Create an empty slice of bytes with the specified length
	bytes := make([]byte, length)
	// Fill the slice with random characters
	for i := range bytes {
		bytes[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(bytes)
}
