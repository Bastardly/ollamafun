package main

import (
	llmhandler "main/llm-handler"
	"net/http"
)

func main() {
	loadUITemplates()
	http.HandleFunc("/generate", llmhandler.HandleGenerate)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	go watchFiles()
	runServer()
}
