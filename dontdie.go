package main

import (
	"fmt"
	"io/ioutil"
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

func createFiles() {
	data := make([]byte, 1024)

	// Start the timer
	startTime := time.Now()

	for i := 0; i < 100000; i++ {
		filename := fmt.Sprintf("file_%d.txt", i)
		err := ioutil.WriteFile(filename, data, os.ModePerm)
		if err != nil {
			fmt.Println("Error writing to file:", filename, err)
		}
	}

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	fmt.Println("10,000 files created successfully!")
	fmt.Printf("Time taken: %s\n", elapsedTime)
}

func main() {
	// Setting up a basic HTTP server
	createFiles()
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
