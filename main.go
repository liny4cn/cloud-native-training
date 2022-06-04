package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
	log.Println("Request URI:", r.RequestURI)
	for key := range r.Header {
		value := r.Header.Get(key)
		w.Header().Add(key, value)
		// 输出 header 日志
		log.Println(key, value)
	}

	// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version := os.Getenv("VERSION")
	w.Header().Add("VERSION", version)

	// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	log.Println("Remote Address:", r.RemoteAddr)
	log.Println("Response Status code:", 200)

	// 4. 当访问 localhost/healthz 时，应返回 200
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")

	// 输出 Response header 日志
	resHeader := w.Header()
	for key := range resHeader {
		log.Println(key, resHeader.Get(key))
	}

}

func logRequest(w http.ResponseWriter, r *http.Request) {
	// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
	log.Println("Request URI:", r.RequestURI)
	for key := range r.Header {
		value := r.Header.Get(key)
		w.Header().Add(key, value)
		// 输出 header 日志
		log.Println(key, value)
	}

	// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version := os.Getenv("VERSION")
	w.Header().Add("VERSION", version)

	// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	log.Println("Remote Address:", r.RemoteAddr)
	log.Println("Response Status code:", 200)

	if strings.HasPrefix(r.RequestURI, "/healthz") {
		// 4. 当访问 localhost/healthz 时，应返回 200
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ok")
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	// 输出 Response header 日志
	resHeader := w.Header()
	for key := range resHeader {
		log.Println(key, resHeader.Get(key))
	}
}

func main() {
	println("Server is running at 0.0.0.0:80.")
	http.HandleFunc("/", logRequest)
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
