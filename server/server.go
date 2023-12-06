package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)


func flashcardsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	openaiApiKey := os.Getenv("OPENAI_API_KEY")

	if openaiApiKey == "" {
		log.Print("OPENAI_API_KEY wasn't provided, some functionality won't work")
	} else {
		aiHandler := AiRequestHandler{openaiApiKey: openaiApiKey}
		http.Handle("/ai", aiHandler)

		ttsServer := makeTtsServer(openaiApiKey)
		http.Handle("/tts", ttsServer)
		defer ttsServer.Close()
	}

	http.HandleFunc("/flashcards", flashcardsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
