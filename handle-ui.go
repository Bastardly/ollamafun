package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

type PageData struct{}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, PageData{})
}

func startUI() {
	var err error
	tpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	http.HandleFunc("/", handleIndex)
}
