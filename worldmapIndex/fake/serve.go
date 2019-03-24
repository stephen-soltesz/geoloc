package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Simple static webserver:
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(os.Getenv("PWD")))))
}
