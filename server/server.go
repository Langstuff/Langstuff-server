package main

import (
	"log"
	"net/http"
	"os"
	"raiden_fumo/lang_notebook_core/database"

	"github.com/joho/godotenv"
)

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

	flashcardHandler := NewFlashcardRequestHandler(database.InitializeDatabase())
	http.HandleFunc("/flashcards/start_session", flashcardHandler.StartSession)
	http.HandleFunc("/flashcards/get_next_card", flashcardHandler.GetNextCard)
	// http.HandleFunc("/flashcards/send_card_answer", flashcardHandler.SendCardAnswer)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
