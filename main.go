package main

import (
	llmhandler "main/llm-handler"
	"net/http"
)

func main() {
	// Hack to access from Ollama on WSL from Windows :s
	// os.Setenv("OLLAMA_HOST", "http://127.0.0.1:8003")

	setMainRoute()
	http.HandleFunc("/generate", llmhandler.HandleGenerate)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	go runServer()
	go watchFiles()

	// Wait indefinitely
	select {}

}
