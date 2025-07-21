package main

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/ollama/ollama/api"
)

var tpl *template.Template

type PageData struct {
	Prompt   string
	Response string
	Error    string
}

func handleUI(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	if r.Method == http.MethodPost {
		r.ParseForm()
		data.Prompt = r.FormValue("prompt")

		client, err := api.ClientFromEnvironment()
		if err != nil {
			data.Error = "Failed to create Ollama client: " + err.Error()
			tpl.Execute(w, data)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		stream := false
		req := &api.GenerateRequest{
			Model:  "llama3.2", // or any model you have pulled
			Prompt: data.Prompt,
			Stream: &stream,
		}

		var response string
		err = client.Generate(ctx, req, api.GenerateResponseFunc(func(res api.GenerateResponse) error {
			response += res.Response
			return nil
		}))
		if err != nil {
			data.Error = "Ollama API error: " + err.Error()
		} else {
			data.Response = response
		}
	}

	tpl.Execute(w, data)
}
