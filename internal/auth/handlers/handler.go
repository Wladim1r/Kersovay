package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"library/internal/auth/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo repository.AuthRepository
}

func NewAuthHandler(repo repository.AuthRepository) *AuthHandler {
	return &AuthHandler{
		repo: repo,
	}
}

func (h *AuthHandler) Register(username, password string) error {
	if len(username) < 3 {
		return fmt.Errorf("❌ имя пользователя должно содержать минимум 3 символа")
	}

	if len(password) < 5 {
		return fmt.Errorf("❌ пароль должен содержать минимум 6 символов")
	}

	exists, err := h.repo.UserExists(username)
	if err != nil {
		return fmt.Errorf("❌ ошибка проверки пользователя: %w", err)
	}

	if exists {
		return fmt.Errorf("❌ пользователь с таким именем уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("❌ ошибка хеширования пароля: %w", err)
	}

	err = h.repo.CreateUser(username, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("❌ пользователь с таким именем уже существует")
		}
		return fmt.Errorf("❌ ошибка создания пользователя: %w", err)
	}

	fmt.Printf("✅ Регистрация успешна! Добро пожаловать, %s!\n\n", username)
	return nil
}

func (h *AuthHandler) Login(username, password string) (int, error) {
	userID, hashedPassword, err := h.repo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("❌ неверное имя пользователя или пароль")
		}
		return 0, fmt.Errorf("❌ ошибка входа: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return 0, fmt.Errorf("❌ неверное имя пользователя или пароль")
	}

	fmt.Printf("✅ Успешный вход! Добро пожаловать, %s!\n\n", username)
	return userID, nil
}
