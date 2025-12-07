package repository

import (
	"database/sql"
	"fmt"
)

type AuthRepository interface {
	CreateUserTable() error
	CreateUser(username, password string) error
	GetUserByUsername(username string) (int, string, error)
	UserExists(username string) (bool, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (ar *authRepository) CreateUserTable() error {
	_, err := ar.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL)`,
	)

	if err != nil {
		return fmt.Errorf("не удалось создать таблицу users: %w", err)
	}
	return nil
}

func (ar *authRepository) CreateUser(username, password string) error {
	_, err := ar.db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		username,
		password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ar *authRepository) GetUserByUsername(username string) (int, string, error) {
	var id int
	var password string

	row := ar.db.QueryRow("SELECT id, password FROM users WHERE username = ?", username)
	err := row.Scan(&id, &password)
	if err != nil {
		return 0, "", err
	}

	return id, password, nil
}

func (ar *authRepository) UserExists(username string) (bool, error) {
	var count int
	row := ar.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
