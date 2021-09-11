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

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Private Chatrooms v1")
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	chatroomId, err := chatroom.Create()

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

	roomError := chatroom.Join(conn, roomId, username)

	if roomError != nil {
		http.Error(w, roomError.Error(), http.StatusBadRequest)
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
