package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// InitDB initializes the connection to the PostgreSQL database and creates tables if they don't exist.
func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Println("DATABASE_URL environment variable not set. Using default local connection string.")
		connStr = "postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	// Create users table
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		nome VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL
	);`

	if _, err = DB.Exec(createUsersTableSQL); err != nil {
		return fmt.Errorf("error creating 'users' table: %w", err)
	}

	// Create pontos table
	createPontosTableSQL := `
	CREATE TABLE IF NOT EXISTS pontos (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		horario TIMESTAMPTZ NOT NULL,
		CONSTRAINT fk_user
			FOREIGN KEY(user_id) 
			REFERENCES users(id)
			ON DELETE CASCADE
	);`

	if _, err = DB.Exec(createPontosTableSQL); err != nil {
		return fmt.Errorf("error creating 'pontos' table: %w", err)
	}

	log.Println("Database initialized and tables are ready.")
	return nil
}