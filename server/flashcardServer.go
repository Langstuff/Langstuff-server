package main

// NOTE: Copilot was used for code generation here

import (
	"encoding/json"
	"net/http"
	"raiden_fumo/lang_notebook_core/core/database"
	"raiden_fumo/lang_notebook_core/core/flashcards"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type flashcardRequestHandler struct {
	db       *gorm.DB
	sessions map[string]*flashcards.Session
}

func NewFlashcardRequestHandler(db *gorm.DB) *flashcardRequestHandler {
	return &flashcardRequestHandler{
		db:       database.InitializeDatabase(),
		sessions: make(map[string]*flashcards.Session, 0),
	}
}

func generateSessionID() string {
	return uuid.New().String()
}

func getFlashcardIDs() []int {
	return []int{}
}


func (handler flashcardRequestHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	session := &flashcards.Session{
		ID:           generateSessionID(),
		CurrentCard:  0,
		FlashcardIDs: getFlashcardIDs(),
	}
	handler.sessions[session.ID] = session
	w.Write([]byte(session.ID))
}

func getSessionIDFromRequest(r *http.Request) string {
	return r.URL.Query().Get("sessionID")
}

func (handler flashcardRequestHandler) GetNextCard(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionIDFromRequest(r)
	session, ok := handler.sessions[sessionID]
	if !ok {
		http.Error(w, "Session "+sessionID+" not found", http.StatusNotFound)
		return
	}
	if session.CurrentCard >= len(session.FlashcardIDs) {
		http.Error(w, "No more cards", http.StatusBadRequest)
		return
	}
	cardID := session.FlashcardIDs[session.CurrentCard]
	var card database.Flashcard
	if err := handler.db.Where("id = ?", cardID).First(&card).Error; err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}
	// TODO: use priority queue or something?
	// Array isn't a way to go here
	session.CurrentCard++
	ret, err := json.Marshal(map[string]interface{}{
		"cardID": cardID,
		"sideA":  card.SideA,
		"sideB":  card.SideB,
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(ret)
}

func (handler flashcardRequestHandler) SendCardAnswer(w http.ResponseWriter, r *http.Request) {
	// sessionID := getSessionIDFromRequest(r)
	// session, ok := handler.sessions[sessionID]
	// if !ok {
	// 	http.Error(w, "Session "+sessionID+" not found", http.StatusNotFound)
	// 	return
	// }
}
