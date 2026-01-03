package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageRequest struct {
	Body string `json:"body" binding:"required"`
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID             uuid.UUID `db:"id" json:"id"`
	ConversationID uuid.UUID `db:"conversation_id" json:"conversation_id"`
	SenderID       uuid.UUID `db:"sender_id" json:"sender_id"`
	Body           string    `db:"body" json:"body"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}