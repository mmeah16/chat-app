package service

import (
	"chat/internal/models"
	"chat/internal/repository"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrSameUserDM   = errors.New("cannot create dm with self")
	ErrUserNotInConversation = errors.New("user is not member of specified conversation")
)

type ConversationService struct {
	DB *sql.DB 
}

func NewConversationService(db *sql.DB) *ConversationService {
	return &ConversationService{
		DB: db,
	}
}

func (s *ConversationService) CreateMessage(ctx context.Context, senderID uuid.UUID, conversationID uuid.UUID, messageRequest *models.MessageRequest) (*models.MessageResponse, error) {
	// Validate conversation exists and sender is apart of the conversation
	valid, err := repository.IsConversationMember(ctx, s.DB, senderID, conversationID)

	if !valid {
		return nil, ErrUserNotInConversation
	}

	message := &models.Message{
		ConversationID: conversationID,
		SenderID: senderID,
		Body: messageRequest.Body,
	}

	messageResponse, err := repository.CreateMessage(ctx, s.DB, message)

	if err != nil {
		return nil, err
	}

	return messageResponse, nil

}

func (s *ConversationService) CreateOrGetDM(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID) (*models.Conversation, error) {
	// Ensure sender and receiver are not the same 
	if senderID == receiverID {
		return nil, ErrSameUserDM
	}

	// Validate receiver ID exists 
	exists, err := repository.UserExists(ctx, s.DB, receiverID)

	if !exists {
		return nil, ErrUserNotFound
	}

	dmKey := computeDMKey(senderID, receiverID)

	// enforces atomicity
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	conv, err := repository.CreateDMConversation(ctx, tx, dmKey)
	if err != nil {
		return nil, err
	}

	if err := repository.AddMember(ctx, tx, conv.ID, senderID); err != nil {
		return nil, err
	}
	if err := repository.AddMember(ctx, tx, conv.ID, receiverID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return conv, nil
}

func computeDMKey(a, b uuid.UUID) string {
	aa, bb := a.String(), b.String()
	if bb < aa {
		aa, bb = bb, aa
	}

	sum := sha256.Sum256([]byte(aa + ":" + bb))
	return hex.EncodeToString(sum[:])
}