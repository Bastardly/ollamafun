package main

import (
	"log"
	"net/http"
	"os"
)

func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6789"
	}
	log.Printf("✅ Server running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
