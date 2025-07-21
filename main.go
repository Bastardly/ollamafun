package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	startUI()
	http.HandleFunc("/generate", handleGenerate)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "6789"
	}
	log.Printf("✅ Server running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
