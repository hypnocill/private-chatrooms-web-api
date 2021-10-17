package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hypnocill/my-private-chatrooms/api/service/chatroom"
)

const (
	ROOM_ID                = "room_id"
	USERNAME               = "username"
	USERNAME_INPUT_TIMEOUT = time.Minute * 5
)

type CreateRoomInput struct {
	Passowrd string `json:"password"`
}

type AuthRoomInput struct {
	RoomId   string `json:"roomId"`
	Passowrd string `json:"password"`
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Private Chatrooms v1")
}

func authRoom(w http.ResponseWriter, r *http.Request) {
	var authRoomInput AuthRoomInput
	password := ""
	roomId := ""

	if err := json.NewDecoder(r.Body).Decode(&authRoomInput); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		password = authRoomInput.Passowrd
		roomId = authRoomInput.RoomId
	}

	authToken, err := chatroom.Auth(roomId, password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseStruct := struct {
		Token string
	}{Token: authToken}

	payload, err := json.Marshal(responseStruct)
	if err != nil {
		http.Error(w, "Could not authenticate", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	var createRoomInput CreateRoomInput
	password := ""
	if err := json.NewDecoder(r.Body).Decode(&createRoomInput); err != nil {
		fmt.Println(err)
	} else {
		password = createRoomInput.Passowrd
	}

	chatroomId, err := chatroom.Create(password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseStruct := struct {
		Id string
	}{Id: chatroomId}

	payload, err := json.Marshal(responseStruct)
	if err != nil {
		http.Error(w, "Could not create chatroom", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error during connection upgrade", http.StatusInternalServerError)
		return
	}

	roomId := mux.Vars(r)[ROOM_ID]    //validate
	username := mux.Vars(r)[USERNAME] //validate

	_, authPayloadJson, _ := conn.ReadMessage()

	var authPayload chatroom.AuthPayload
	_ = json.Unmarshal(authPayloadJson, &authPayload)
	hashedPassword := authPayload.Message

	roomError := chatroom.Join(conn, roomId, username, hashedPassword)

	if roomError != nil {
		fmt.Println(roomError.Error())
		// http.Error(w, roomError.Error(), http.StatusBadRequest)
		cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, roomError.Error())
		conn.WriteMessage(websocket.CloseMessage, cm)
		time.Sleep(time.Second * 2)
		conn.Close()

		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { //remove
		return true
	},
}
