package handler

import (
	"encoding/json"
	"fmt"
	"go-conversation/model"
	"go-conversation/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SendChatHandler handles sending messages between users.
func SendChatHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Memeriksa keberadaan token di header Authorization
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		// Memverifikasi token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Mendapatkan informasi pengguna dari token
		userID, ok := claims["id"].(string)
		if !ok {
			http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
			return
		}

		var newMessage struct {
			Sender   string `json:"sender"`
			Receiver string `json:"receiver"`
			Message  string `json:"message"`
		}

		err = json.NewDecoder(r.Body).Decode(&newMessage)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if newMessage.Sender == "" || newMessage.Receiver == "" || newMessage.Message == "" {
			http.Error(w, "Sender, Receiver, and Message are required fields", http.StatusBadRequest)
			return
		}

		// Memastikan bahwa pengirim pesan sesuai dengan informasi dari token
		if newMessage.Sender != userID {
			http.Error(w, "Unauthorized: Sender does not match token claims", http.StatusUnauthorized)
			return
		}

		// Create a Conversation object with auto-generated values
		conversation := model.Conversation{
			Sender:   newMessage.Sender,
			Receiver: newMessage.Receiver,
			Message:  newMessage.Message,
			DateTime: time.Now(),
		}

		// Insert the new message into the conversations collection
		conversationID, err := repo.AddMessage(r.Context(), conversation)
		if err != nil {
			fmt.Println("Error adding message:", err)
			http.Error(w, "Error adding message", http.StatusInternalServerError)
			return
		}

		// Mengisi respons yang diinginkan
		response := model.CustomDataResponse{
			StatusCode: http.StatusCreated,
			Message:    "Pesan berhasil dikirim",
			Data: model.UserData{
				ConversationID: conversationID,
				Sender:         conversation.Sender,
				Receiver:       conversation.Receiver,
				Message:        conversation.Message,
				DateTime:       conversation.DateTime.UTC().String(),
			},
		}

		// Mengirim respons dalam format JSON
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// GetUserConversationHandler handles retrieving all conversations for a user with optional time filters.
func GetUserConversationHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Memeriksa keberadaan token di header Authorization
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		// Memverifikasi token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		/// Mendapatkan informasi pengguna dari token
        userID, ok := claims["id"].(string)
        if !ok {
            http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
            return
        }

		// Mendapatkan parameter query string
		queryValues := r.URL.Query()
		afterUnixtimeStr := queryValues.Get("after_unixtime")
		beforeUnixtimeStr := queryValues.Get("before_unixtime")
		limitStr := queryValues.Get("limit")

		var afterUnixtime, beforeUnixtime int64
		var limit int

		// Konversi string menjadi int64
		if afterUnixtimeStr != "" {
			afterUnixtime, err = strconv.ParseInt(afterUnixtimeStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid after_unixtime parameter", http.StatusBadRequest)
				return
			}
		}

		if beforeUnixtimeStr != "" {
			beforeUnixtime, err = strconv.ParseInt(beforeUnixtimeStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid before_unixtime parameter", http.StatusBadRequest)
				return
			}
		}

		// Konversi string menjadi int
		if limitStr != "" {
    		limit, err = strconv.Atoi(limitStr)
    		if err != nil {
        		http.Error(w, "Parameter limit tidak valid: harus berupa bilangan bulat positif yang valid", http.StatusBadRequest)
        		return
    		}

    		// Validasi limit untuk memastikan nilainya tidak negatif
    		if limit < 0 {
        		http.Error(w, "Parameter limit tidak valid: harus berupa bilangan bulat positif yang valid", http.StatusBadRequest)
        		return
    		}
		}
		// Mengambil semua percakapan untuk pengguna berdasarkan user_id dan parameter waktu
		conversations, err := repo.GetUserConversationsWithTimeFilter(r.Context(), userID, afterUnixtime, beforeUnixtime, limit)
		if err != nil {
			fmt.Println("Error getting user conversations:", err)
			http.Error(w, "Error getting user conversations", http.StatusInternalServerError)
			return
		}

		// Mengisi respons yang diinginkan
		response := model.CustomDataResponse{
			StatusCode: http.StatusOK,
			Message:    "List of user conversations",
			Data: model.UserData{
				Username:       claims["username"].(string),
				ID:             userID,
				ConversationID: "",
				Conversations:  conversations,
			},
		}

		// Mengirim respons dalam format JSON
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
