package main

import (
	"io"
	"log"
	"net/http"
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
	log.Println("Remote Address:", r.RemoteAddr)

	// 4. 当访问 localhost/healthz 时，应返回 200
	if r.RequestURI == "/healthz" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ok")
		log.Println("Response Status code:", 200)
	} else {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Response Status code:", 404)
	}

	// 输出 Response header 日志
	resHeader := w.Header()
	for key := range resHeader {
		log.Println(key, resHeader.Get(key))
	}
}

func main() {
	println("Server is running at 0.0.0.0:80.")
	mux := http.NewServeMux()
	mux.HandleFunc("/", logRequest)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatalf("Server start failed: %v", err.Error())
		return
	}
}
