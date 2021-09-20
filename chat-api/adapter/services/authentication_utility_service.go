package services

import (
	"chat-api/adapter/logger"
	"chat-api/domain"
	"chat-api/infrastructure/config"
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationUtility struct {
	log logger.Logger
}

func NewAuthenticationUtility(log logger.Logger) AuthenticationUtility {
	return AuthenticationUtility{
		log: log,
	}
}

func (a AuthenticationUtility) HashPassword(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil

}

func (a AuthenticationUtility) CheckPasswordHash(ctx context.Context, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}

func (a AuthenticationUtility) GenerateToken(ctx context.Context, user domain.User) (string, error) {
	cfg := config.GetConfig()

	// var err error
	// atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	// atClaims["first_name"] = user.FirstName()
	// atClaims["last_name"] = user.LastName()
	// atClaims["email"] = user.Email()
	// atClaims["role"] = user.Role()
	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	// at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// token, err := at.SignedString([]byte(cfg.AccessSecret))
	// if err != nil {
	// 	return "", err
	// }
	// return token, nil
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["first_name"] = user.FirstName()
	claims["last_name"] = user.LastName()
	claims["email"] = user.Email()
	claims["role"] = user.Role()
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(cfg.AccessSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
