package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// >>>>> added for alloy

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Alloy tails this JSON log file while stdout remains useful locally.
	log.Logger = zerolog.New(io.MultiWriter(os.Stdout, logFile)).With().Timestamp().Logger()

	// <<<< added for alloy

	log.Debug().
		Str("Scale", "833 cents").
		Str("app", "go-app").
		Float64("Interval", 833.09).
		Msg("Fibonacci is everywhere")

	log.Debug().
		Str("app", "go-app").
		Str("Name", "Tom").
		Send()

	log.Info().
		Str("app", "go-app").
		Str("request_id", "abc123").
		Msg("order created")
}

// Output: {"level":"debug","Scale":"833 cents","Interval":833.09,"time":1562212768,"message":"Fibonacci is everywhere"}
// Output: {"level":"debug","Name":"Tom","time":1562212768}
