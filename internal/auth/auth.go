package auth

import (
	"errors"
	"horseshoe-server/internal/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var JwtSecret []byte

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func Register(username, password string) error {
	hash, err := CreateHash(password)
	if err != nil {
		return err
	}

	user := db.User{
		Username:     username,
		PasswordHash: hash,
	}

	result := db.DB.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func Login(username, password string) (string, string, error) {
	var user db.User

	result := db.DB.Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", "", ErrUserNotFound
	} else if result.Error != nil {
		return "", "", result.Error
	}

	match, err := ComparePassword(password, user.PasswordHash)
	if err != nil || !match {
		return "", "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, user.Username, nil
}
