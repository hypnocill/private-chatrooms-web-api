package chatroom

import "encoding/json"

type AuthPayload struct {
	PayloadType PAYLOAD_TYPE
	Username    string
	Message     string
}

func NewAuthPayload(username string, message string) *AuthPayload {
	payload := AuthPayload{Username: username, Message: message, PayloadType: PAYLOAD_TYPE_AUTH}

	return &payload
}

func (authPayload *AuthPayload) toJson() ([]byte, error) {
	return json.Marshal(authPayload)
}
