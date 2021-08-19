package chatroom

import "encoding/json"

const (
	MESSAGE_CATEGORY_TEXT = iota + 1
	MESSAGE_CATEGORY_USERNAME
)

type MessagePayload struct {
	Username string
	Message  string
	Category int
}

func NewMessagePayload(username string, message string, category int) *MessagePayload {
	payload := MessagePayload{Username: username, Message: string(message), Category: MESSAGE_CATEGORY_TEXT}

	return &payload
}

func (this *MessagePayload) toJson() ([]byte, error) {
	return json.Marshal(this)
}
