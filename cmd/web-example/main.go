package main

import (
	"log"
	"net/http"
)

func main() {
	// Set the directory that contains your HTML and WASM files
	fs := http.FileServer(http.Dir("./public"))

	// Handle all requests by serving a file of the same name
	http.Handle("/", fs)

	// Start the server on port 8080
	log.Println("Listening on http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
