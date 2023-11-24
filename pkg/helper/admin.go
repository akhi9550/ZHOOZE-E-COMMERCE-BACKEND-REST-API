package helper

import (
	"Zhooze/pkg/config"
	"Zhooze/pkg/utils/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type authCustomClaimsAdmin struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	jwt.StandardClaims
}

func GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, error) {
	cfg, _ := config.LoadConfig()
	claims := &authCustomClaimsAdmin{
		Firstname: admin.Firstname,
		Lastname:  admin.Lastname,
		Email:     admin.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.KEY_ADMIN))
	if err != nil {
		fmt.Println("Error is", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*authCustomClaimsAdmin, error) {
	cfg, _ := config.LoadConfig()
	token, err := jwt.ParseWithClaims(tokenString, &authCustomClaimsAdmin{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.KEY_ADMIN), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*authCustomClaimsAdmin); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
