package models

type User struct {
	ID       int    `json:"id"`       // Уникальный идентификатор пользователя
	Email    string `json:"email"`    // Email пользователя
	Password string `json:"password"` // Пароль пользователя
}
