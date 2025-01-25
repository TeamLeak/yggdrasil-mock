package models

import "time"

type Texture struct {
	ID         int       `json:"id"`          // Уникальный идентификатор текстуры
	Hash       string    `json:"hash"`        // Хэш текстуры
	Data       []byte    `json:"data"`        // Данные текстуры
	UploadedAt time.Time `json:"uploaded_at"` // Дата и время загрузки текстуры
}
