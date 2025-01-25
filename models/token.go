package models

import "time"

type Token struct {
	ID          int       `json:"id"`           // Уникальный идентификатор токена
	AccessToken string    `json:"access_token"` // Токен доступа
	ClientToken string    `json:"client_token"` // Токен клиента
	CreatedAt   time.Time `json:"created_at"`   // Время создания токена
	UserID      int       `json:"user_id"`      // Идентификатор пользователя
	CharacterID int       `json:"character_id"` // Идентификатор персонажа (может быть NULL)
}
