package repository

import (
	"chat/internal/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func CreateUser(db *sql.DB, user *models.User) error {

	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	return db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	var user models.User

	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserExists(ctx context.Context, db *sql.DB, userID uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	if err := db.QueryRowContext(ctx, query, userID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}