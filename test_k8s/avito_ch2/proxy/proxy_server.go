package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	timeServerURL, ok := os.LookupEnv("ADDR")
	if !ok {
		slog.Error("missing env variable ADDR")
		os.Exit(1)
	}

	getTimeURL := fmt.Sprintf("%s/time", timeServerURL)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("received request to /hello from " + r.RemoteAddr)

		resp, err := http.DefaultClient.Get(getTimeURL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Time server (%s) is not available: %v", timeServerURL, err)
			return
		}
		defer resp.Body.Close()

		respTime, _ := io.ReadAll(resp.Body)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello World at %s", string(respTime))
	})

	slog.Info("starting proxy server on port 8890")

	err := http.ListenAndServe(":8890", nil)
	if err != nil {
		slog.Error("failed to start proxy server:", err)
	}
}
