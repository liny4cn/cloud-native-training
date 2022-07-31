package main

import (
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func metricsRequest(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("URI", r.RequestURI).Msg("Enter Request")
	timer := NewTimer()
	defer timer.ObserveTotal()

	delay := rand.Intn(2000)
	time.Sleep(time.Duration(delay+10) * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello world.")

	log.Info().Msg("Response end.")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("URI", r.RequestURI).Msg("Enter Request")
	timer := NewTimer()
	defer timer.ObserveTotal()

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")

	log.Info().Int("Response Status code", http.StatusOK).Msg("Response end.")
}
