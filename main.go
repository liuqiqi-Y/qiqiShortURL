package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Trace 用于记录日志
var Trace *log.Logger

func init() {
	file, err := os.OpenFile("trace.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file: ", err)
	}
	Trace = log.New(io.MultiWriter(file, os.Stderr), "Trace: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	router := mux.NewRouter()
	InitializeRouters(router)
	Run(":8000", router)
}

// Run run
func Run(address string, router *mux.Router) {
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}
	server.ListenAndServe()
}

// InitializeRouters 初始化路由
func InitializeRouters(router *mux.Router) {
	router.HandleFunc("/api/shortURL", CreateShortURL).Methods("POST")
	router.HandleFunc("/api/shortURLInfo", ShortURLInfo).Methods("GET")
	router.HandleFunc("/{shortURL:[a-zA-Z0-9]{1,11}}", Redirect).Methods("GET")
}
