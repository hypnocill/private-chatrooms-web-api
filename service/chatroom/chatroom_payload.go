package chatroom

import "encoding/json"

type ChatroomPayload struct {
	Id    string
	Topic string
}

func NewChatroomPayload(id string, topic string) *ChatroomPayload {
	payload := ChatroomPayload{Id: id, Topic: topic}

	return &payload
}

func (this *ChatroomPayload) toJson() ([]byte, error) {
	return json.Marshal(this)
}
