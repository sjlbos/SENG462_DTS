package main

import (
    "log"
    "net/http"
    "time"
)


func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        println("Starting Request")
        start := time.Now()

        inner.ServeHTTP(w, r)

        log.Printf(
		"%10s%20s%20s",
		r.Header.Get("X-TransNo"),
		name,
		time.Since(start),
        )
        println("Finished Request\n")
    })
}
