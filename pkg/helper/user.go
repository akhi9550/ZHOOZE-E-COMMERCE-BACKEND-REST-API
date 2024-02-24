package helper

import (
	"Zhooze/pkg/config"
	"Zhooze/pkg/utils/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func PasswordHash(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("internal server error")
	}
	hash := string(hashPassword)
	return hash, nil
}

func GenerateTokenUsers(userID int, userEmail string, expirationTime time.Time) (string, error) {
	cfg, _ := config.LoadConfig()
	claims := &AuthUserClaims{
		Id:    userID,
		Email: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateAccessToken(user models.UserDetailsResponse) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	tokenString, err := GenerateTokenUsers(user.Id, user.Email, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(user models.UserDetailsResponse) (string, error) {
	expirationTime := time.Now().Add(24 * 90 * time.Hour)
	tokenString, err := GenerateTokenUsers(user.Id, user.Email, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func CompareHashAndPassword(a string, b string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a), []byte(b))
	if err != nil {
		return err
	}
	return nil
}
func PasswordHashing(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("internal server error")
	}
	hash := string(hashedPassword)
	return hash, nil
}
