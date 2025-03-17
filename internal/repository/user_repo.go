package repository

import (
	"context"
	"log"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password, role string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New()

	query := `INSERT INTO users (id, username, password, role) VALUES ($1, $2, $3, $4) RETURNING created_at`
	row := database.DB.QueryRow(context.Background(), query, id, username, string(hashedPassword), role)

	var createdAt string
	err = row.Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	log.Println("User created successfully")
	return &models.User{ID: id, Username: username, Role: role}, nil
}

func GetUser(username string) (*models.User, error) {
	query := `SELECT id, username, password, role, created_at FROM users WHERE username=$1`
	row := database.DB.QueryRow(context.Background(), query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
