package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // O _ significa que estamos importando pelo efeito colateral (registrar o driver).
)

var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados e cria a tabela se ela não existir.
func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o banco de dados: %w", err)
	}

	// Tenta criar a tabela. Se ela já existir, o comando não fará nada.
	createTableSQL := `CREATE TABLE IF NOT EXISTS pontos (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"horario" DATETIME NOT NULL
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("erro ao criar a tabela 'pontos': %w", err)
	}

	log.Println("Banco de dados inicializado e tabela 'pontos' pronta.")
	return nil
}