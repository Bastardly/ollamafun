package llmhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ollama/ollama/api"
)

type generateInput struct {
	Prompt string `json:"prompt"`
}

func (g generateInput) getInsultJSON() string {
	return `
You are a content moderation AI. Analyze the following input and respond only with a JSON object. 

Mark the content as offensive if it contains:
  - Hate speech, slurs, or threats
  - Harassment or demeaning language directed at specific people
  - Praise or justification of harmful ideologies (e.g. racism, slavery, fascism, Nazism)
  - Praise of individuals widely recognized for crimes against humanity (e.g. Hitler)
  - Explicit insults or abusive language targeted directly at a person (e.g. "I think you are an idiot")

Do NOT mark the content as offensive if it contains:
  - Negative opinions or emotional reactions about non-human entities like:
    - Games, products, books, movies, software, food, brands, or everyday objects
    - Example: "I hate Star Citizen", "I hate Cyberpunk 2077", or "This game is trash"
  - Criticism or negative opinions about public figures or celebrities
  - Hate or speech against individuals known for crimes against humanity (e.g. Hitler)

The JSON must have:
  - "offensive": a boolean indicating if the text is offensive,
  - "reason": a short explanation if it is offensive (or null if not offensive).

Input: """` + g.Prompt + `"""
Respond in JSON only.
`
}

func (g generateInput) getCharacter() string {
	return `
Du er en AI, der skriver farverige og sjove karakterintroduktioner. Analyser følgende input og svar kun med et JSON-objekt.

Hvis prompten indeholder en karakter, så:
  - Brug alle detaljer fra prompten.
  - Sørg for at modtageren ved, hvem karakteren er, hvad de hedder, hvad de kan lide, og hvordan de opfører sig.
  - Undgå klichéer – find på noget originalt og levende.
  - Gør introduktionen som om karakteren selv præsenterer sig fx. ("Hej, jeg hedder...").
  - Sørg for at det føles personligt, sjovt og lidt skævt.
  - Tilføj gerne små finurlige detaljer, vaner eller særheder.

Hvis prompten **ikke** indeholder en karakter, så opfind en selv ud fra din fantasi. Gør den mindeværdig og underholdende.

Svar skal være på **dansk**, og ikke på svensk eller norsk. Det skal være i følgende JSON-format:

{
  "navn": "Karakterens fulde navn på dansk",
  "introduktion": "En kort og personlig introduktion skrevet i jeg-form på dansk",
  "beskrivelse": "En detaljeret beskrivelse i 3. person af karakterens personlighed, baggrund og særpræg på dansk"
}

Input: """` + g.Prompt + `"""
Svar KUN med JSON – ingen forklaringer eller ekstra tekst.
`
}

func HandleGenerate(w http.ResponseWriter, r *http.Request) {
	var input generateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		http.Error(w, "failed to create Ollama client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	req := &api.GenerateRequest{
		Model:  "llama3.2",
		Prompt: input.getCharacter(),
		Stream: func(b bool) *bool { return &b }(false),
		Options: map[string]interface{}{
			"temperature": 0.4,
			"max_tokens":  800,
		},
	}

	var result string
	err = client.Generate(ctx, req, api.GenerateResponseFunc(func(res api.GenerateResponse) error {
		result += res.Response
		return nil
	}))
	if err != nil {
		http.Error(w, "generation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"response": result,
	})
}
