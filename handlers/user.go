package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jscastaneda-esp/rest-ws-go/models"
	"github.com/jscastaneda-esp/rest-ws-go/repository"
	"github.com/jscastaneda-esp/rest-ws-go/server"
	"github.com/jscastaneda-esp/rest-ws-go/services"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	SignUpRequest
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := new(SignUpRequest)
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			log.Println("Json Decode:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			log.Println("NewRandom:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			log.Println("GenerateFromPassword:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := &models.User{
			BaseModel: models.BaseModel{
				Id: id.String(),
			},
			Email:    request.Email,
			Password: string(hashedPassword),
		}
		err = repository.InsertUser(r.Context(), user)
		if err != nil {
			log.Println("InsertUser:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := new(LoginRequest)
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			log.Println("Json Decode:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			log.Println("GetUserByEmail:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			log.Println("Not found user:", err)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			log.Println("CompareHashAndPassword:", err)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		claims := &models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			log.Println("SignedString:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})
	}
}

func MeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, status, err := services.GetClaimsToken(r.Header.Get("Authorization"), s.Config().JWTSecret)
		if err != nil {
			log.Println("GetUserData:", err)
			http.Error(w, err.Error(), status)
			return
		}

		user, err := repository.GetUserById(r.Context(), claims.UserId)
		if err != nil {
			log.Println("GetUserById:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
