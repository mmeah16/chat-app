package repository

import (
	"chat/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateMessage(ctx context.Context, db *sql.DB, message *models.Message) (*models.MessageResponse, error)  {
	query := `
		INSERT INTO messages (conversation_id, sender_id, body)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`
	var id uuid.UUID
	var createdAt time.Time

	err := db.QueryRowContext(
		ctx,
		query,
		message.ConversationID,
		message.SenderID,
		message.Body,
	).Scan(&id, &createdAt)
	
	if err != nil {
		return nil, fmt.Errorf("CreateMessage: %w", err)
	}

	return &models.MessageResponse{
		ID:        id,
		SenderID:  message.SenderID, 
		Body:      message.Body,     
		CreatedAt: createdAt,
	}, nil
}

func GetMessagesByConversation(ctx context.Context, db *sql.DB, conversationID uuid.UUID) ([]*models.MessageResponse, error) {
	query := `
		SELECT id, sender_id, body, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := db.QueryContext(ctx, query, conversationID)

	if err != nil {
		return nil, fmt.Errorf("GetMessagesByConversation: %w", err)
	}

	defer rows.Close()

	var messages []*models.MessageResponse

	for rows.Next() {
		var msg models.MessageResponse

		if err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.Body,
			&msg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan message row: %w", err)
		}

		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
} 