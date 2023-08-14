package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Printf("Started server with pid %d\n", os.Getpid())

	signal.Ignore(syscall.SIGTERM)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	instance, ok := os.LookupEnv("CF_INSTANCE_INDEX")
	if !ok {
		instance = "UNKNOWN"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = writer.Write([]byte(fmt.Sprintf("This could be a meaningful HTTP response coming from instance %s\n", instance)))
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
