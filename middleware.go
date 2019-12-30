package main

import (
	"net/http"
	"time"
)

// LogMiddleWare 记录一个请求消耗的时间
func LogMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(t1)
		Trace.Println(timeElapsed)
	})
}

// RecoverHandler 恢复panic
func RecoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Trace.Println("Recover from panic: ", err)
				http.Error(w, "server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
