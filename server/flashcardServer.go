package main

// NOTE: Copilot was used for code generation here

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type flashcardRequestHandler struct {
	openaiApiKey string
}

func generateSessionID() string {
	return uuid.New().String()
}

func getFlashcardIDs() []int {
	return []int{}
}

func (handler flashcardRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}

type Session struct {
	ID           string
	CurrentCard  int
	FlashcardIDs []int
}

var sessions = make(map[string]*Session)

func (handler flashcardRequestHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	session := &Session{
		ID:           generateSessionID(),
		CurrentCard:  0,
		FlashcardIDs: getFlashcardIDs(),
	}
	sessions[session.ID] = session
	// Return the session ID to the client
}

func getSessionIDFromRequest(r *http.Request) string {
	return r.URL.Query().Get("sessionID")
}

func (handler flashcardRequestHandler) GetNextCard(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	sessionID := getSessionIDFromRequest(r)
	session, ok := sessions[sessionID]
	if !ok {
		// Handle error: session not found
		return
	}
	if session.CurrentCard >= len(session.FlashcardIDs) {
		// Handle error: no more cards
		return
	}
	cardID := session.FlashcardIDs[session.CurrentCard]
	// Retrieve the card from the database using the cardID
	var card Flashcard
	if err := db.Where("id = ?", cardID).First(&card).Error; err != nil {
		// Handle error: card not found
		http.Error(w, "Card not found", http.StatusNotFound)
		return nil, err
	}
	session.CurrentCard++
	ret, err := json.Marshal(map[string]int{"cardID": cardID, "sideA": card.SideA, "sideB": card.SideB})
	return ret, err
}
