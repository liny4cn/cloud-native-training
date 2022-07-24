package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogger(envLevel string) {
	println("Init logger with ", envLevel)

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
		println("Invalid log level:", envLevel, ", Current log level is INFO.")
		log.Warn().Str("log-level", envLevel).Msgf("Current log level is INFO.")
		return
	}

	println("Current log level is", envLevel, "(", logLevel, ") .")
}
