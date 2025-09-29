package services

import (
	"database/sql"

	"github.com/DevaSinha/StreamSight/go-api/config"
	"github.com/DevaSinha/StreamSight/go-api/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(email, password string) (*models.User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = config.DB.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at",
		email, string(hashedPwd),
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	return &user, err
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, email, password FROM users WHERE email=$1", email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, nil // User not found
	}
	return &user, err
}

func ValidateCredentials(email, password string) (*models.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, nil // Invalid password
	}

	return user, nil
}
