package database

import (
	"database/sql"
	"time"
	"yggdrasil/models"
)

func GetUserByID(db *sql.DB, userID int) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`SELECT id, email, password FROM users WHERE id = ?`, userID).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmailOrCharacter(db *sql.DB, identifier string) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT u.id, u.email, u.password 
		FROM users u 
		LEFT JOIN characters c ON u.id = c.user_id 
		WHERE u.email = ? OR c.name = ?`, identifier, identifier).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindCharacterByName(db *sql.DB, name string) (*models.Character, error) {
	var character models.Character
	err := db.QueryRow(`SELECT id, uuid, name, model, user_id FROM characters WHERE name = ?`, name).
		Scan(&character.ID, &character.UUID, &character.Name, &character.Model, &character.UserID)
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func RevokeAllTokens(db *sql.DB, userID int) error {
	_, err := db.Exec(`DELETE FROM tokens WHERE user_id = ?`, userID)
	return err
}

func DeleteToken(db *sql.DB, accessToken string) error {
	_, err := db.Exec(`DELETE FROM tokens WHERE access_token = ?`, accessToken)
	return err
}

func GetToken(db *sql.DB, accessToken string) (*models.Token, error) {
	var token models.Token
	err := db.QueryRow(`SELECT id, access_token, client_token, created_at, user_id, character_id FROM tokens WHERE access_token = ?`, accessToken).
		Scan(&token.ID, &token.AccessToken, &token.ClientToken, &token.CreatedAt, &token.UserID, &token.CharacterID)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func InsertToken(db *sql.DB, token *models.Token) error {
	_, err := db.Exec(`INSERT INTO tokens (access_token, client_token, created_at, user_id, character_id) VALUES (?, ?, ?, ?, ?)`,
		token.AccessToken, token.ClientToken, token.CreatedAt, token.UserID, token.CharacterID)
	return err
}

func GetUserCharacters(db *sql.DB, userID int) []map[string]string {
	rows, err := db.Query(`SELECT uuid, name FROM characters WHERE user_id = ?`, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var characters []map[string]string
	for rows.Next() {
		var uuid, name string
		if err := rows.Scan(&uuid, &name); err == nil {
			characters = append(characters, map[string]string{
				"id":   uuid,
				"name": name,
			})
		}
	}
	return characters
}

func GetCharacterByUUID(db *sql.DB, uuid string) (*models.Character, error) {
	var character models.Character
	err := db.QueryRow(
		`SELECT id, uuid, name, model, user_id FROM characters WHERE uuid = ?`,
		uuid,
	).Scan(&character.ID, &character.UUID, &character.Name, &character.Model, &character.UserID)
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func GetCharactersByNames(db *sql.DB, names []string) ([]map[string]string, error) {
	var profiles []map[string]string

	query := `SELECT uuid, name FROM characters WHERE name = ?`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, name := range names {
		var uuid, charName string
		err := stmt.QueryRow(name).Scan(&uuid, &charName)
		if err == nil {
			profiles = append(profiles, map[string]string{
				"id":   uuid,
				"name": charName,
			})
		}
	}

	return profiles, nil
}

func FindCharacterByServerAndName(db *sql.DB, name, serverID string) (*models.Character, error) {
	var character models.Character
	err := db.QueryRow(
		`SELECT c.id, c.uuid, c.name, c.model, c.user_id
		 FROM characters c
		 JOIN tokens t ON c.id = t.character_id
		 WHERE c.name = ? AND t.client_token = ?`,
		name, serverID,
	).Scan(&character.ID, &character.UUID, &character.Name, &character.Model, &character.UserID)
	if err != nil {
		return nil, err
	}
	return &character, nil
}

// InsertOrUpdateTexture inserts or updates a texture in the database
func InsertOrUpdateTexture(db *sql.DB, uuid, textureType string, data []byte) error {
	hash := uuid + "-" + textureType
	_, err := db.Exec(
		`INSERT INTO textures (hash, data, uploaded_at) 
		 VALUES (?, ?, ?)
		 ON CONFLICT(hash) DO UPDATE 
		 SET data = excluded.data, uploaded_at = excluded.uploaded_at`,
		hash, data, time.Now(),
	)
	return err
}

// DeleteTexture deletes a texture from the database by UUID and type
func DeleteTexture(db *sql.DB, uuid, textureType string) error {
	hash := uuid + "-" + textureType
	_, err := db.Exec(`DELETE FROM textures WHERE hash = ?`, hash)
	return err
}

// GetTextureByHash retrieves a texture by its hash
func GetTextureByHash(db *sql.DB, hash string) (*struct {
	Hash       string
	Data       []byte
	UploadedAt time.Time
}, error) {
	var texture struct {
		Hash       string
		Data       []byte
		UploadedAt time.Time
	}
	err := db.QueryRow(
		`SELECT hash, data, uploaded_at FROM textures WHERE hash = ?`,
		hash,
	).Scan(&texture.Hash, &texture.Data, &texture.UploadedAt)
	if err != nil {
		return nil, err
	}
	return &texture, nil
}

func GetAllUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query(`SELECT id, email, password FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func InsertUser(db *sql.DB, user *models.User) error {
	_, err := db.Exec(`INSERT INTO users (email, password) VALUES (?, ?)`, user.Email, user.Password)
	return err
}

func UpdateUser(db *sql.DB, id string, user *models.User) error {
	_, err := db.Exec(`UPDATE users SET email = ?, password = ? WHERE id = ?`, user.Email, user.Password, id)
	return err
}

func DeleteUser(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

// SetCharacterTexture - Set a texture for a character
func SetCharacterTexture(db *sql.DB, uuid, textureHash string) error {
	_, err := db.Exec(`UPDATE characters SET texture_hash = ? WHERE uuid = ?`, textureHash, uuid)
	return err
}

// InsertCharacter - Add a new character to the database
func InsertCharacter(db *sql.DB, character *models.Character) error {
	_, err := db.Exec(
		`INSERT INTO characters (uuid, name, model, user_id) VALUES (?, ?, ?, ?)`,
		character.UUID, character.Name, character.Model, character.UserID,
	)
	return err
}
