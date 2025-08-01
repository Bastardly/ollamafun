package main

import (
	llmhandler "main/llm-handler"
	"net/http"
)

func main() {
	setMainRoute()
	http.HandleFunc("/generate", llmhandler.Chat)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	go runServer()

	// Development reload tool for client watch
	go watchFiles()

	// Wait indefinitely
	select {}

}
