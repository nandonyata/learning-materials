package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

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
