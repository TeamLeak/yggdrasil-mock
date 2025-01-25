package database

import (
	"database/sql"
	"fmt"
	"log"
	"yggdrasil/config"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
	_ "modernc.org/sqlite"             // SQLite driver
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	var dsn string
	switch cfg.Database.Type {
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name,
		)
	case "sqlite":
		dsn = cfg.Database.SQLiteFile
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
		)
	default:
		log.Printf("Unsupported database type: %s", cfg.Database.Type)
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	db, err := sql.Open(cfg.Database.Type, dsn)
	if err != nil {
		log.Printf("Failed to open the database connection: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("Failed to verify the database connection: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully.")
	return db, nil
}

func Migrate(db *sql.DB, cfg *config.Config) error {
	log.Println("Starting database migration...")

	var schema string
	switch cfg.Database.Type {
	case "sqlite":
		schema = `
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
	case "postgres":
		schema = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS characters (
	id SERIAL PRIMARY KEY,
	uuid TEXT UNIQUE NOT NULL,
	name TEXT UNIQUE NOT NULL,
	model TEXT NOT NULL DEFAULT 'STEVE',
	user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tokens (
	id SERIAL PRIMARY KEY,
	access_token TEXT UNIQUE NOT NULL,
	client_token TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	user_id INTEGER NOT NULL REFERENCES users (id),
	character_id INTEGER REFERENCES characters (id)
);

CREATE TABLE IF NOT EXISTS textures (
	id SERIAL PRIMARY KEY,
	hash TEXT UNIQUE NOT NULL,
	data BYTEA NOT NULL,
	uploaded_at TIMESTAMP NOT NULL
);
`
	case "mysql":
		schema = `
CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS characters (
	id INT AUTO_INCREMENT PRIMARY KEY,
	uuid VARCHAR(36) UNIQUE NOT NULL,
	name VARCHAR(255) UNIQUE NOT NULL,
	model VARCHAR(255) NOT NULL DEFAULT 'STEVE',
	user_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tokens (
	id INT AUTO_INCREMENT PRIMARY KEY,
	access_token VARCHAR(255) UNIQUE NOT NULL,
	client_token VARCHAR(255) NOT NULL,
	created_at DATETIME NOT NULL,
	user_id INT NOT NULL,
	character_id INT,
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (character_id) REFERENCES characters (id)
);

CREATE TABLE IF NOT EXISTS textures (
	id INT AUTO_INCREMENT PRIMARY KEY,
	hash VARCHAR(255) UNIQUE NOT NULL,
	data BLOB NOT NULL,
	uploaded_at DATETIME NOT NULL
);
`
	default:
		log.Printf("Unsupported database type: %s", cfg.Database.Type)
		return fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("Failed to execute schema: %v", err)
		return err
	}

	log.Println("Database migration completed successfully.")
	return nil
}
