package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ROUTE_INDEX       = "/"
	ROUTE_ROOM_JOIN   = "/room/join/%s"
	ROUTE_ROOM_CREATE = "/room/create"
)

func SetupRoutes() {
	router := mux.NewRouter()

	router.HandleFunc(ROUTE_INDEX, index).Methods("GET")
	router.HandleFunc(fmt.Sprintf(ROUTE_ROOM_JOIN, "{room_id}"), joinRoom).Methods("GET")
	router.HandleFunc(ROUTE_ROOM_CREATE, createRoom).Methods("GET")

	http.Handle(ROUTE_INDEX, router)
}
