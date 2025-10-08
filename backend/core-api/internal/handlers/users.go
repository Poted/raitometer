package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Poted/raitometer/backend/core-api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

func (h *Handlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad Request: could not decode JSON", http.StatusBadRequest)
		return
	}

	if input.Email == "" || len(input.Password) < 8 {
		http.Error(w, "Bad Request: email is required and password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email: input.Email,
	}

	err = user.SetPassword(input.Password)
	if err != nil {
		log.Printf("error setting password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.userStore.Create(user)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) && pqError.Code == "23505" {
			http.Error(w, "Conflict: user with this email already exists", http.StatusConflict)
		} else {
			log.Printf("error creating user: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	err = h.writeJSON(w, http.StatusCreated, user, nil)
	if err != nil {
		log.Printf("error writing register user response: %v", err)
	}
}

func (h *Handlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad Request: could not decode JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userStore.GetByEmail(input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized: invalid credentials", http.StatusUnauthorized)
		} else {
			log.Printf("error getting user by email: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	passwordMatch, err := user.CheckPassword(input.Password)
	if err != nil {
		log.Printf("error checking password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !passwordMatch {
		http.Error(w, "Unauthorized: invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "raitometer-api",
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := claims.SignedString(jwtSecret)
	if err != nil {
		log.Printf("error signing token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	err = h.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		log.Printf("error writing login response: %v", err)
	}
}
