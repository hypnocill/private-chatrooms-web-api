package chatroom

type PAYLOAD_TYPE int

const (
	PAYLOAD_TYPE_CHATROOM PAYLOAD_TYPE = iota + 1
	PAYLOAD_TYPE_MESSAGE
)
