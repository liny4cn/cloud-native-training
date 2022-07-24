package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/gobike/envflag"
	"github.com/rs/zerolog/log"
)

func startHttpServer(bindPort string, quit chan error) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/", logRequest)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{Addr: bindPort, Handler: mux}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			quit <- err
		}
	}()

	return srv
}

func main() {
	var (
		bindPort string
		serverAddr string
		envLevel string
	)

	flag.StringVar(&serverAddr, "addr", "0.0.0.0", "Server address")
	flag.StringVar(&bindPort, "port", "80", "Server port")
	flag.StringVar(&envLevel, "log-level", "INFO", "logger level")
	envflag.Parse()

	initLogger(envLevel)

	fmt.Println("Server version v1.0.0.")
	fmt.Printf("Server is running at %s:%s.\n", serverAddr, bindPort)

	quit := make(chan error, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	httpSrv := startHttpServer(serverAddr + ":" + bindPort, quit)

	for {
		select {
		case err := <-quit:
			if err == http.ErrServerClosed {
				log.Info().Msg("Server closed.")
			} else {
				log.Error().Err(err).Msg("Server start failed.")
			}
			fmt.Println("Server closed.")
			return
		case s := <-sigs:
			fmt.Printf("Recived notify signal %s (%d), Server is going to shudown.\n", s.String(), s)
			log.Warn().Str("signal", s.String()).Msg("Server is going to shutdown.")
			httpSrv.Shutdown(context.Background())
		}
	}
}
