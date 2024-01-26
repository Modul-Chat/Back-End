package model

import "time"

// Conversation represents the conversation data model.
type Conversation struct {
	ID        string    `json:"conversation_id" bson:"_id,omitempty"`
	Sender    string    `json:"sender" bson:"sender"`
	Receiver  string    `json:"receiver" bson:"receiver"`
	Message   string    `json:"message" bson:"message"`
	DateTime  time.Time `json:"dateTime" bson:"dateTime"`
}

// ConversationsResponse struct untuk format respons daftar percakapan pengguna
type ConversationsResponse struct {
	StatusCode int           `json:"status_code"`
	Message    string        `json:"message"`
	Data       []Conversation `json:"data"`
}
