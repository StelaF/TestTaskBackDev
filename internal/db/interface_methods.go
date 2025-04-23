package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"time"
)

func (db *DataBase) Get(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
	statement := `SELECT token_hash, expires_at 
				  FROM refresh_tokens 
				  WHERE user_id = $1;`

	var tokenHash string
	var expiresAt time.Time

	err := db.conn.QueryRowContext(ctx, statement, userID).Scan(&tokenHash, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, err
		}
		return "", time.Time{}, err
	}

	return tokenHash, expiresAt, nil
}

func (db *DataBase) Store(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	statement := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) ON CONFLICT (user_id)DO UPDATE SET token_hash = $2, expires_at = $3;`

	_, err := db.conn.ExecContext(ctx, statement, userID, tokenHash, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) Delete(ctx context.Context, userID uuid.UUID) error {
	statement := `DELETE FROM refresh_tokens 
       			  WHERE user_id = $1;`

	_, err := db.conn.ExecContext(ctx, statement, userID)
	if err != nil {
		return err
	}

	return nil
}
