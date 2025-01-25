package models

type Character struct {
	ID          int    `json:"id"`      // Уникальный идентификатор персонажа
	UUID        string `json:"uuid"`    // UUID персонажа
	Name        string `json:"name"`    // Имя персонажа
	Model       string `json:"model"`   // Модель персонажа (например, "STEVE")
	UserID      int    `json:"user_id"` // Идентификатор пользователя-владельца персонажа
	TextureHash string `json:"texture_hash"`
}
