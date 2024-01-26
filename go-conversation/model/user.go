package model

import "time"

// User represents the user data model.
type User struct {
	ID        string    `json:"user_id" bson:"_id,omitempty"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

// Response struct untuk format respons yang diinginkan
type Response struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

// UserData struct untuk menyimpan data pengguna dalam respons
type UserData struct {
	Username        string `json:"username,omitempty"`
	CreatedAt       string `json:"createdAt,omitempty"`
	ID              string `json:"id,omitempty"`
	ConversationID  string `json:"conversation_id,omitempty"`
	Sender          string `json:"sender,omitempty"`
	Receiver        string `json:"receiver,omitempty"`
	Message         string `json:"message,omitempty"`
	DateTime        string `json:"dateTime,omitempty"`
	Token           string `json:"token,omitempty"`
	Conversations   []Conversation `json:"conversations,omitempty"`
}

