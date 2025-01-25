package services

import (
	"database/sql"
	_ "yggdrasil/utils"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUserByEmailOrCharacter(db *sql.DB, identifier string) (*User, error) {
	var user User
	err := db.QueryRow(
		`SELECT u.id, u.email, u.password
		 FROM users u
		 LEFT JOIN characters c ON u.id = c.user_id
		 WHERE u.email = ? OR c.name = ?`,
		identifier, identifier,
	).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
