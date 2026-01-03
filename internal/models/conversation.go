package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ConversationType string

const (
	ConversationTypeDM    ConversationType = "dm"
	ConversationTypeGroup ConversationType = "group"
)

type Conversation struct {
	ID        uuid.UUID        `db:"id" json:"id"`
	Type      ConversationType `db:"type" json:"type"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
	DmKey     sql.NullString   `json:"-"`
}

type ConversationMembers struct {
	ConversationID uuid.UUID `db:"conversation_id" json:"conversation_id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	JoinedAt       time.Time `db:"joined_at" json:"joined_at"`
}


func (t ConversationType) IsValid() bool {
	switch t {
	case ConversationTypeDM, ConversationTypeGroup:
		return true
	default:
		return false
	}
}