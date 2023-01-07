package main

import (
	"context"
	"expvar"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	// LogStrKeyModule is for use with the logger as a key to specify the module name.
	LogStrKeyModule = "module"
	// LogStrKeyService is for use with the logger as a key to specify the service name.
	LogStrKeyService = "service"
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
	// Start API Service

	log.Info().Msg("main: Initializing Service support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	api := http.Server{}

	// Start the service listening for requests.
	go func() {
		log.Info().Msgf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	ctx := log.WithContext(context.Background())

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "service error")

	case sig := <-shutdown:
		log.Info().Msgf("main: %v: Start shutdown", sig)

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			if err := api.Close(); err != nil {
				return errors.Wrap(err, "could not stop server gracefully")
			}
			return errors.Wrap(err, "could not stop server")
		}
	}

	return nil
}
