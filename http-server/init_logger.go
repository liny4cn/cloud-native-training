package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogger(envLevel string) {
	fmt.Printf("Init logger with %s.\n", envLevel)

	logLevel := zerolog.InfoLevel
	isMatched := true
	switch envLevel {
	case "PANIC":
		logLevel = zerolog.PanicLevel
	case "FATAL":
		logLevel = zerolog.FatalLevel
	case "ERROR":
		logLevel = zerolog.ErrorLevel
	case "WARN":
		logLevel = zerolog.WarnLevel
	case "INFO":
		logLevel = zerolog.InfoLevel
	case "DEBUG":
		logLevel = zerolog.DebugLevel
	case "TRACE":
		logLevel = zerolog.TraceLevel
	default:
		isMatched = false
	}
	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	if !isMatched {
		fmt.Printf("Invalid log level: %s, Current log level is INFO.", envLevel)
		log.Warn().Str("log-level", envLevel).Msg("Current log level is INFO.")
		return
	}

	fmt.Printf("Current log level is %s(%d).\n", envLevel, logLevel)
}
