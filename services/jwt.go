package services

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/jscastaneda-esp/rest-ws-go/models"
)

func CheckToken(tokenString string, jwtSecret string) (*jwt.Token, error) {
	tokenString = strings.TrimSpace(tokenString)
	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, errors.New("invalid credentials")
	}

	return jwt.ParseWithClaims(tokenParts[1], &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}

func GetClaimsToken(tokenString string, jwtSecret string) (*models.AppClaims, int, error) {
	token, err := CheckToken(tokenString, jwtSecret)
	if err != nil {
		log.Println("CheckToken:", err)
		return nil, http.StatusUnauthorized, err
	}

	if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		return claims, 0, nil
	} else {
		return nil, http.StatusUnauthorized, errors.New("token invalid")
	}
}
