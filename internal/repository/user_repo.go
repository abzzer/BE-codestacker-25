package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user models.UserCreate) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	query := `
		INSERT INTO users (name, password, role, clearance_level)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id string
	err = database.DB.QueryRow(context.Background(), query, user.Name, string(hashedPassword), string(user.Role), string(user.ClearanceLevel)).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func DeleteUser(id string) error {
	query := `UPDATE users SET deleted = TRUE WHERE id = $1;`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}

func UpdateUser(input models.UpdateUser) error {
	if input.ID == "" {
		return errors.New("user ID is required")
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *input.Name)
		argIdx++
	}

	if input.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, fmt.Sprintf("password = $%d", argIdx))
		args = append(args, string(hashed))
		argIdx++
	}

	if input.Role != nil {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, *input.Role)
		argIdx++
	}

	if input.ClearanceLevel != nil {
		setClauses = append(setClauses, fmt.Sprintf("clearance_level = $%d", argIdx))
		args = append(args, *input.ClearanceLevel)
		argIdx++
	}

	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	args = append(args, input.ID)

	query := fmt.Sprintf(`
		UPDATE users SET %s
		WHERE id = $%d AND deleted = FALSE;
	`, strings.Join(setClauses, ", "), argIdx)

	_, err := database.DB.Exec(context.Background(), query, args...)
	return err
}

func GetPasswordAndRoleByUserID(id string) (string, string, error) {
	query := `
		SELECT password, role
		FROM users
		WHERE id = $1 AND deleted = FALSE;
	`

	var password, role string
	err := database.DB.QueryRow(context.Background(), query, id).Scan(&password, &role)
	if err != nil {
		return "", "", errors.New("invalid user ID or user is deleted")
	}

	return password, role, nil
}
