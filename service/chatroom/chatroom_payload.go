package chatroom

import "encoding/json"

type ChatroomPayload struct {
	PayloadType PAYLOAD_TYPE
	Id          string
	Topic       string
	Users       map[string]struct{}
}

func NewChatroomPayload(id string, topic string, users map[string]struct{}, payloadType PAYLOAD_TYPE) *ChatroomPayload {
	payload := ChatroomPayload{Id: id, Topic: topic, Users: users, PayloadType: payloadType}

	return &payload
}

func (this *ChatroomPayload) toJson() ([]byte, error) {
	return json.Marshal(this)
}
