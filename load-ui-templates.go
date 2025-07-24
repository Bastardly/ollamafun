package main

import (
	"html/template"
	"log"
	llmhandler "main/llm-handler"
	"net/http"
	"sync"
	"time"
)

var templates *template.Template
var changedTime time.Time

var (
	mutex = &sync.Mutex{}
)

func updateUITemplates() {
	mutex.Lock()
	defer mutex.Unlock()
	log.Println("Clearing template cache...")

	templates = template.New("clear")
	changedTime = time.Now()
	loadUITemplates()
}

type PageData struct {
	PromptMethods []string
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	templates.Execute(w, PageData{
		PromptMethods: llmhandler.AvailableMethods,
	})
}

func loadUITemplates() {
	var err error
	templates, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
}

func setMainRoute() {
	http.HandleFunc("/", handleIndex)
}
