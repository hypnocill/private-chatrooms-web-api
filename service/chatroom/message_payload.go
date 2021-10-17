package chatroom

import "encoding/json"

type MessagePayload struct {
	PayloadType PAYLOAD_TYPE
	Username    string
	Message     string
}

func NewMessagePayload(username string, message string) *MessagePayload {
	payload := MessagePayload{Username: username, Message: message, PayloadType: PAYLOAD_TYPE_MESSAGE}

	return &payload
}

func (messagePayload *MessagePayload) toJson() ([]byte, error) {
	return json.Marshal(messagePayload)
}
