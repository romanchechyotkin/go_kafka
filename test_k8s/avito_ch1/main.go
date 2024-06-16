package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("received request to /hello from " + r.RemoteAddr)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello World")
	})

	slog.Info("starting server on port 8890")
	err := http.ListenAndServe(":8890", nil)
	if err != nil {
		slog.Error("error starting server: ", err)
	}
}
