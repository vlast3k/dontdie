package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	instance, ok := os.LookupEnv("CF_INSTANCE_INDEX")
	if !ok {
		instance = "UNKNOWN"
	}
	_, _ = w.Write([]byte(fmt.Sprintf("This could be a meaningful HTTP response coming from instance %s\n", instance)))
}

func main() {
	// Setting up a basic HTTP server
	http.HandleFunc("/", handler)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	// Setting up signal catching
	sigCh := make(chan os.Signal, 1)

	// Catching SIGTERM
	signal.Notify(sigCh, syscall.SIGTERM)

	// Waiting for SIGTERM
	sig := <-sigCh
	if sig == syscall.SIGTERM {
		fmt.Println("SIGTERM received. Waiting for 5 minutes...")

		// Waiting for 5 minutes
		time.Sleep(5 * time.Minute)
		fmt.Println("5 minutes elapsed, shutting down.")
	}
}
