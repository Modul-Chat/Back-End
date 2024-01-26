package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-conversation/model"
	"go-conversation/repository"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles user registration.
func RegisterHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser model.User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Check if the required fields are empty
		if newUser.Username == "" || newUser.Password == "" {
			http.Error(w, "Username and Password are required fields", http.StatusBadRequest)
			return
		}

		// Check if the username is already taken
		existingUser, err := repo.GetUserByUsername(r.Context(), newUser.Username)
		if err == nil && existingUser.ID != "" {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}

		// Insert the new user into the database
		userID, err := repo.CreateUser(r.Context(), newUser)
		if err != nil {
			fmt.Println("Error creating user:", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		newUser.ID = userID

		// Mengisi respons yang diinginkan (tanpa menyertakan token)
		response := model.CustomDataResponse{
			StatusCode: http.StatusCreated,
			Message:    "Register berhasil",
			Data: model.UserData{
				Username:  newUser.Username,
				CreatedAt: time.Now().UTC().String(),
				ID:        newUser.ID,
			},
		}

		// Mengirim respons dalam format JSON
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// LoginHandler handles user login.
func LoginHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Retrieve user from the database
		user, err := repo.GetUserByUsername(r.Context(), credentials.Username)
		if err != nil {
			http.Error(w, "Invalid username", http.StatusUnauthorized)
			return
		}

		// Compare hashed password with the provided password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Omit password from the response
		user.Password = ""
		user.CreatedAt = time.Time{}

		// Generate JWT token
		token, err := generateJWTToken(user, "your-secret-key")
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Mengisi respons yang diinginkan
		response := model.CustomDataResponse{
			StatusCode: http.StatusOK,
			Message:    "Login berhasil",
			Data: model.UserData{
				Username: user.Username,
				ID:       user.ID,
				Token:    token,
			},
		}

		// Mengirim respons dalam format JSON
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}