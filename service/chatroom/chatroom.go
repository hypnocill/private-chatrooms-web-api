package chatroom

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/hypnocill/my-private-chatrooms/api/service/cache"
)

var ctx = context.Background()

const DEFAULT_EXPIRATION_TIME = time.Hour

func Auth(roomId string, password string) (string, error) {
	rdb := cache.OpenRedisConnection()

	info := rdb.Ping(ctx)
	if info.Err() != nil {
		return "", errors.New("could not establish connection")
	}

	chatroomPayloadJson := rdb.Get(ctx, roomId)
	if chatroomPayloadJson.Err() == redis.Nil {
		return "", errors.New("room does not exist")
	}

	var chatroomPayload ChatroomPayload
	chatroomPayloadBytes, _ := chatroomPayloadJson.Bytes()
	_ = json.Unmarshal(chatroomPayloadBytes, &chatroomPayload)
	fmt.Println("payload", chatroomPayload.Passowrd)
	hashedPassword := sha256.Sum256([]byte(password))
	encodedHashedPassword := hex.EncodeToString(hashedPassword[:])
	fmt.Println("incoming password", encodedHashedPassword)

	if chatroomPayload.Passowrd != encodedHashedPassword {
		return "", errors.New("password is wrong")
	}

	return encodedHashedPassword, nil
}

func Create(password string) (string, error) {
	rdb := cache.OpenRedisConnection()

	info := rdb.Ping(ctx)
	if info.Err() != nil {
		return "", errors.New("could not establish connection")
	}

	// add salt
	hashedPassword := sha256.Sum256([]byte(password))
	encodedHashedPassword := hex.EncodeToString(hashedPassword[:])

	roomId := makeRandomString()
	emptyUsers := make(map[string]struct{})

	chatroomPayload := NewChatroomPayload(roomId, "", emptyUsers, encodedHashedPassword)
	chatroomPayloadJson, _ := chatroomPayload.toJson() //handle error
	rdb.Set(ctx, roomId, chatroomPayloadJson, DEFAULT_EXPIRATION_TIME)

	return roomId, nil
}

func Join(conn *websocket.Conn, roomId string, username string, hashedPassword string) error {
	rdbPublisher := cache.OpenRedisConnection()

	info := rdbPublisher.Ping(ctx)
	if info.Err() != nil {
		return errors.New("could not establish connection")
	}

	chatroomPayloadJson := rdbPublisher.Get(ctx, roomId)
	if chatroomPayloadJson.Err() == redis.Nil {
		return errors.New("room does not exist")
	}

	// Move this adding of a User to chatroom to 'chatroom_payload'
	var chatroomPayload ChatroomPayload
	chatroomPayloadBytes, _ := chatroomPayloadJson.Bytes()
	_ = json.Unmarshal(chatroomPayloadBytes, &chatroomPayload)
	fmt.Println("payload", chatroomPayload.Passowrd)
	fmt.Println("incoming password", hashedPassword)

	if chatroomPayload.Passowrd != hashedPassword {
		return errors.New("password is wrong")
	}

	_, userExists := chatroomPayload.Users[username]

	if userExists {
		return errors.New("username is taken")
	}

	var Empty struct{}
	chatroomPayload.Users[username] = Empty

	newChatroomPayloadJson, _ := chatroomPayload.toJson()

	rdbPublisher.Set(ctx, roomId, newChatroomPayloadJson, DEFAULT_EXPIRATION_TIME)

	rdbSubscriber := cache.OpenRedisConnection().Subscribe(ctx, roomId)

	conn.SetReadDeadline(time.Now().Add(5 * time.Minute)) //adjust

	go streamMessagesFromUserToChatroomChannel(conn, rdbPublisher, roomId, username)
	go streamMessagesFromChatroomChannelToUser(conn, rdbSubscriber, 1)

	notifyChatroomForNewUser(rdbPublisher, chatroomPayload, roomId)

	return nil
}

func notifyChatroomForNewUser(rdbPublisher *redis.Client, chatroomPayload ChatroomPayload, roomId string) {
	chatroomPayload.Passowrd = "***"
	newChatroomPayloadJson, _ := chatroomPayload.toJson()
	rdbPublisher.Publish(ctx, roomId, newChatroomPayloadJson)
}

func streamMessagesFromUserToChatroomChannel(conn *websocket.Conn, rdbPublisher *redis.Client, roomId string, username string) {
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			rdbPublisher.Close()
			// remove user from chatroom payload here
			log.Println("Error during message reading:", err)
			break
		}

		payload := NewMessagePayload(username, string(message))
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
