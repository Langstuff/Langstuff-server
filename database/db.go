package database

// AI Note: Copilot was used for code generation here

import (
	"gorm.io/gorm"
)

type Deck struct {
	gorm.Model
	ID          uint `gorm:"unique,primaryKey,autoIncrement:true"`
	Name        string
	Description string
	Flashcards  []Flashcard
}

type Flashcard struct {
	gorm.Model
	ID        uint `gorm:"unique,primaryKey,autoIncrement:true"`
	SideA     string
	SideB     string
	DeckID    uint      `gorm:"foreignKey:DeckID"`
	SM2Record SM2Record `gorm:"foreignKey:FlashcardID"`
	Tags      []Tag     `gorm:"many2many:flashcard_tag_pairs;"`
}

type SM2Record struct {
	gorm.Model
	FlashcardID    uint `gorm:"primaryKey;foreignKey:FlashcardID"`
	EasinessFactor float64
	Interval       uint
	Repetition     uint
}

type Tag struct {
	gorm.Model
	ID   uint `gorm:"unique,primaryKey,autoIncrement:true"`
	Name string
}

func (tag *Tag) GetFlashcards(db *gorm.DB) []Flashcard {
	var flashcards []Flashcard

	// Query the FlashcardTagPair table to get flashcard IDs associated with the tag
	var pairs []FlashcardTagPair
	db.Where("tag_id = ?", tag.ID).Find(&pairs)

	// Get the flashcard IDs
	var flashcardIDs []uint
	for _, pair := range pairs {
		flashcardIDs = append(flashcardIDs, pair.FlashcardID)
	}

	// Query the Flashcard table to retrieve the flashcards with the matching IDs
	db.Where("id IN ?", flashcardIDs).Find(&flashcards)

	return flashcards
}

func GetFlashcardsByTagList(db *gorm.DB, deckName string, tags []string) []Flashcard {
	if deckName == "" {
		return getFlashcardsByTags(db, tags)
	} else {
		return getFlashcardsByDeckAndTags(db, deckName, tags)
	}
}

func getFlashcardsByTags(db *gorm.DB, tags []string) []Flashcard {
	var flashcards []Flashcard
	db.Preload("Tags").
		Joins(`JOIN flashcard_tag_pairs
			ON flashcards.id = flashcard_tag_pairs.flashcard_id`).
		Joins("JOIN tags ON flashcard_tag_pairs.tag_id = tags.id").
		Where("tags.name IN ?", tags).
		Find(&flashcards)
	return flashcards
}

func getFlashcardsByDeckAndTags(db *gorm.DB, deckName string, tags []string) []Flashcard {
	var deck Deck
	db.Where("name = ?", deckName).First(&deck)

	var flashcards []Flashcard
	db.Preload("Tags").
		Joins(`JOIN flashcard_tag_pairs
			ON flashcards.id = flashcard_tag_pairs.flashcard_id
			AND flashcards.deck_id = ?`, deck.ID).
		Joins("JOIN tags ON flashcard_tag_pairs.tag_id = tags.id").
		Where("tags.name IN ?", tags).
		Find(&flashcards)
	return flashcards
}

type FlashcardTagPair struct {
	gorm.Model
	FlashcardID uint `gorm:"primaryKey;foreignKey:FlashcardID"`
	TagID       uint `gorm:"primaryKey;foreignKey:TagID"`
}
