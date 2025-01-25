package services

import (
	"database/sql"
	"time"
)

type Token struct {
	ID          int       `json:"id"`
	AccessToken string    `json:"access_token"`
	ClientToken string    `json:"client_token"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`
	CharacterID int       `json:"character_id"`
}

func InsertToken(db *sql.DB, token *Token) error {
	_, err := db.Exec(
		`INSERT INTO tokens (access_token, client_token, created_at, user_id, character_id) VALUES (?, ?, ?, ?, ?)`,
		token.AccessToken, token.ClientToken, token.CreatedAt, token.UserID, token.CharacterID,
	)
	return err
}

func GetToken(db *sql.DB, accessToken string) (*Token, error) {
	var token Token
	err := db.QueryRow(
		`SELECT id, access_token, client_token, created_at, user_id, character_id FROM tokens WHERE access_token = ?`,
		accessToken,
	).Scan(&token.ID, &token.AccessToken, &token.ClientToken, &token.CreatedAt, &token.UserID, &token.CharacterID)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func DeleteToken(db *sql.DB, accessToken string) error {
	_, err := db.Exec(`DELETE FROM tokens WHERE access_token = ?`, accessToken)
	return err
}

func RevokeAllTokens(db *sql.DB, userID int) error {
	_, err := db.Exec(`DELETE FROM tokens WHERE user_id = ?`, userID)
	return err
}
