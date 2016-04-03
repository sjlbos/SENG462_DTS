package main

import (
    "log"
    "net/http"
    "time"
    
)


func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        inner.ServeHTTP(w, r)

        log.Printf(
		"%10s%20s%20s",
		r.Header.Get("X-TransNo"),
		name,
		time.Since(start),
        )
    })
}
