package services

import (
	"context"
	"database/sql"
	"federal-funds-rate-metrics-ByYear/db"
	"federal-funds-rate-metrics-ByYear/dto"
	"federal-funds-rate-metrics-ByYear/models"
	"fmt"
)

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]models.User, error) {
	query := "SELECT id, name, email FROM users"
	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return users, nil
}

// GetUserByID retrieves a user by ID from the database
func GetUserByID(id int) (models.User, error) {
	query := "SELECT id, name, email FROM users WHERE id = $1"
	var user models.User

	err := db.Conn.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil // Return an empty user if not found
		}
		return models.User{}, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return user, nil
}

// CreateUser adds a new user to the database and returns the created user
func CreateUser(requestDto *dto.UserDto) (*dto.UserDto, error) {
	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	var newUserID int

	err := db.Conn.QueryRow(context.Background(), query, requestDto.Name, requestDto.Email).Scan(&newUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &dto.UserDto{
		Name:  requestDto.Name,
		Email: requestDto.Email,
	}, nil
}
