package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func logRequest(w http.ResponseWriter, r *http.Request) {
	// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
	log.Info().Str("URI", r.RequestURI).Msg("Request")
	for key, value := range r.Header {
		w.Header().Add(key, strings.Join(value, ","))
		// 输出 header 日志
		log.Info().Str(key, strings.Join(value, ",")).Msg("Request")
	}

	// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version := os.Getenv("VERSION")
	if version == "" {
		w.Header().Add("VERSION", "1.0")
	} else {
		w.Header().Add("VERSION", version)
	}

	// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	remoteIP := getRealIP(r)
	w.Header().Add("X-Real-Ip", remoteIP)
	log.Info().Str("Remote Address", remoteIP).Msg("Request")

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello world.")
	log.Info().Int("Response Status code", http.StatusOK).Msg("Response")

	// 输出 Response header 日志
	resHeader := w.Header()
	for key := range resHeader {
		log.Info().Str(key, resHeader.Get(key)).Msg("Response")
	}
}

// 4. 当访问 localhost/healthz 时，应返回 200
func healthz(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("URI", r.RequestURI).Msg("Request")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")
	log.Info().Int("Response Status code", http.StatusOK).Msg("Response")
}

func getRealIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0]); ip != "" {
		return ip
	}

	if ip := strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
