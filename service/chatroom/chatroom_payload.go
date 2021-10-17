package chatroom

import "encoding/json"

type ChatroomPayload struct {
	PayloadType PAYLOAD_TYPE
	Id          string
	Topic       string
	Users       map[string]struct{}
	Passowrd    string
}

func NewChatroomPayload(id string, topic string, users map[string]struct{}, password string) *ChatroomPayload {
	payload := ChatroomPayload{Id: id, Topic: topic, Users: users, PayloadType: PAYLOAD_TYPE_CHATROOM, Passowrd: password}

	return &payload
}

func (this *ChatroomPayload) toJson() ([]byte, error) {
	return json.Marshal(this)
}
