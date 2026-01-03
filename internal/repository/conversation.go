package repository

import (
	"chat/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateConversation(ctx context.Context, db *sql.DB, convType models.ConversationType) (*models.Conversation, error) {
	query := `
		INSERT INTO conversations (type)
		VALUES ($1)
		RETURNING id, type, created_at;
	`

	var conv models.Conversation

	err := db.QueryRowContext(
		ctx,
		query,
		convType,
	).Scan(&conv.ID, &conv.Type, &conv.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("CreateConversation: %w", err)
	}

	return &conv, nil
}

func GetConversationsByUser(ctx context.Context, db *sql.DB, userID uuid.UUID) ([]*models.Conversation, error) {

	query := `
		SELECT c.id, c.type, c.created_at
		FROM conversations c
		JOIN conversation_members cm
			ON cm.conversation_id = c.id
		WHERE cm.user_id = $1
		ORDER BY c.created_at DESC;
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetConversationsByUser: %w", err)
	}
	defer rows.Close()

	var conversations []*models.Conversation

	for rows.Next() {
		var conv models.Conversation

		if err := rows.Scan(
			&conv.ID,
			&conv.Type,
			&conv.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan conversation: %w", err)
		}

		conversations = append(conversations, &conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return conversations, nil
}

func CreateDMConversation(ctx context.Context, tx *sql.Tx, dmKey string) (*models.Conversation, error) {
	const insertQuery = `
		INSERT INTO conversations (type, dm_key)
		VALUES ('dm', $1)
		ON CONFLICT (dm_key) WHERE type = 'dm'
		DO NOTHING
		RETURNING id, type, created_at;
	`

	var conv models.Conversation
	var createdAt time.Time

	row := tx.QueryRowContext(ctx, insertQuery, dmKey)
	err := row.Scan(&conv.ID, &conv.Type, &createdAt)


	switch {
	// Conversation record was created
	case err == nil:
		conv.CreatedAt = createdAt
		conv.DmKey = sql.NullString{String: dmKey, Valid: true}
		return &conv, nil
	
	// Insert skipped due to existing DM key (ON CONFLICT DO NOTHING → RETURNING no rows)
	case errors.Is(err, sql.ErrNoRows):
		const selectQ = `
			SELECT id, type, created_at
			FROM conversations
			WHERE type = 'dm' AND dm_key = $1
			LIMIT 1;
		`

		if err := tx.QueryRowContext(ctx, selectQ, dmKey).Scan(&conv.ID, &conv.Type, &conv.CreatedAt); err != nil {
			return nil, fmt.Errorf("select existing dm conversation: %w", err)
		}
		conv.DmKey = sql.NullString{String: dmKey, Valid: true}
		return &conv, nil

	default: 
		return nil, fmt.Errorf("insert dm conversation: %w", err)
	}
}

func AddMember(ctx context.Context, tx *sql.Tx, conversationID, userID uuid.UUID) error {
	const q = `
		INSERT INTO conversation_members (conversation_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (conversation_id, user_id) DO NOTHING;
	`
	if _, err := tx.ExecContext(ctx, q, conversationID, userID); err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func IsConversationMember(ctx context.Context, db *sql.DB, senderID, conversationID uuid.UUID) (bool, error) {
	const query = `
		SELECT EXISTS
		(SELECT cm.user_id 
		FROM conversation_members cm 
		WHERE conversation_id = $1 
		AND user_id = $2)
	`
	var exists bool 

	if err := db.QueryRowContext(ctx, query, conversationID, senderID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}