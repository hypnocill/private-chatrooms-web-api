package chatroom

import "encoding/json"

type MessagePayload struct {
	Username string
	Message  string
}

func NewMessagePayload(username string, message string) *MessagePayload {
	payload := MessagePayload{Username: username, Message: string(message)}

	return &payload
}

func (messagePayload *MessagePayload) toJson() ([]byte, error) {
	return json.Marshal(messagePayload)
}
