package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS characters (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uuid TEXT UNIQUE NOT NULL,
	name TEXT UNIQUE NOT NULL,
	model TEXT NOT NULL DEFAULT 'STEVE',
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tokens (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	access_token TEXT UNIQUE NOT NULL,
	client_token TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	user_id INTEGER NOT NULL,
	character_id INTEGER,
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (character_id) REFERENCES characters (id)
);

CREATE TABLE IF NOT EXISTS textures (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	hash TEXT UNIQUE NOT NULL,
	data BLOB NOT NULL,
	uploaded_at DATETIME NOT NULL
);
`

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "yggdrasil_mock.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("Failed to execute schema: %v", err)
		return err
	}
	return nil
}
