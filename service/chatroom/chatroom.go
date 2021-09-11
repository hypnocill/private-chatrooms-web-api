package chatroom

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/hypnocill/my-private-chatrooms/api/service/cache"
)

var ctx = context.Background()

const DEFAULT_EXPIRATION_TIME = time.Hour

func Create() (string, error) {
	rdb := cache.OpenRedisConnection()

	info := rdb.Ping(ctx)
	if info.Err() != nil {
		return "", errors.New("could not establish connection")
	}

	roomId := makeRandomString()

	chatroomPayload := NewChatroomPayload(roomId, "", make(map[string]struct{}), PAYLOAD_TYPE_CHATROOM)
	chatroomPayloadJson, _ := chatroomPayload.toJson() //handle error
	rdb.Set(ctx, roomId, chatroomPayloadJson, DEFAULT_EXPIRATION_TIME)

	return roomId, nil
}

func Join(conn *websocket.Conn, roomId string, username string) error {
	rdbPublisher := cache.OpenRedisConnection()

	info := rdbPublisher.Ping(ctx)
	if info.Err() != nil {
		return errors.New("could not establish connection")
	}

	chatroomPayloadJson := rdbPublisher.Get(ctx, roomId)
	if chatroomPayloadJson.Err() == redis.Nil {
		return errors.New("room does not exist")
	}

	var chatroomPayload ChatroomPayload
	chatroomPayloadBytes, _ := chatroomPayloadJson.Bytes()
	_ = json.Unmarshal(chatroomPayloadBytes, &chatroomPayload)

	_, userExists := chatroomPayload.Users[username]

	if userExists {
		return errors.New("username is taken")
	}

	var Empty struct{}
	chatroomPayload.Users[username] = Empty
	newChatroomPayloadJson, _ := chatroomPayload.toJson()

	rdbPublisher.Set(ctx, roomId, newChatroomPayloadJson, DEFAULT_EXPIRATION_TIME)
	rdbPublisher.Publish(ctx, roomId, newChatroomPayloadJson)

	rdbSubscriber := cache.OpenRedisConnection().Subscribe(ctx, roomId)

	go streamMessagesFromUserToChatroomChannel(conn, rdbPublisher, roomId, username)
	go streamMessagesFromChatroomChannelToUser(conn, rdbSubscriber, 1)

	return nil
}

func streamMessagesFromUserToChatroomChannel(conn *websocket.Conn, rdbPublisher *redis.Client, roomId string, username string) {
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			rdbPublisher.Close()
			log.Println("Error during message reading:", err)
			break
		}

		payload := NewMessagePayload(username, string(message), PAYLOAD_TYPE_MESSAGE)
		payloadJson, _ := payload.toJson()
		rdbPublisher.Publish(ctx, roomId, payloadJson)
	}
}

func streamMessagesFromChatroomChannelToUser(conn *websocket.Conn, rdbSubscriber *redis.PubSub, messageType int) {
	defer conn.Close()
	chatChannel := rdbSubscriber.Channel()

	for redisMessage := range chatChannel {
		err := conn.WriteMessage(messageType, []byte(redisMessage.Payload))
		if err != nil {
			rdbSubscriber.Close()
			log.Println("Error during message writing:", err)
			break
		}
	}
}

func makeRandomString() string { //@todo refactor to return a safely random string, not timestamp string
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}
