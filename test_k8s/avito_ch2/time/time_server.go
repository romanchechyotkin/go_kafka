package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("received request to /hello from " + r.RemoteAddr)

		t := time.Now()

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%d", t.Unix())
	})

	slog.Info("starting time server on :8891")

	err := http.ListenAndServe(":8891", nil)
	if err != nil {
		slog.Error("failed to start time server:", err)
	}
}
