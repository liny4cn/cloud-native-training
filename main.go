package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

func logRequest(w http.ResponseWriter, r *http.Request) {
	// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
	log.Println("Request URI:", r.RequestURI)
	for key, value := range r.Header {
		w.Header().Add(key, strings.Join(value, ","))
		// 输出 header 日志
		log.Println(key, ":", value)
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
	log.Println("Remote Address:", remoteIP)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello world.")
	log.Println("Response Status code:", 200)

	// 输出 Response header 日志
	resHeader := w.Header()
	for key := range resHeader {
		log.Println(key, ":", resHeader.Get(key))
	}
}

// 4. 当访问 localhost/healthz 时，应返回 200
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")
	log.Println("Response Status code:", 200)
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

func main() {
	println("Server is running at 0.0.0.0:80.")
	mux := http.NewServeMux()
	mux.HandleFunc("/", logRequest)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatalf("Server start failed: %v", err.Error())
		return
	}
}
