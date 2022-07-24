package main

import (
	"flag"
	"net/http"
	"net/http/pprof"

	"github.com/gobike/envflag"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		bindPort string
		envLevel string
	)

	flag.StringVar(&bindPort, "bind-port", ":80", "Server port")
	flag.StringVar(&envLevel, "log-level", "INFO", "Server port")
	envflag.Parse()

	initLogger(envLevel)

	println("Server is running at ", bindPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", logRequest)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err := http.ListenAndServe(bindPort, mux)
	if err != nil {
		log.Error().Err(err).Msg("Server start failed.")
		return
	}
}
