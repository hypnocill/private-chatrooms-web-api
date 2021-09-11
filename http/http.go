package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	ROUTE_INDEX       = "/"
	ROUTE_ROOM_JOIN   = "/room/join/%s/%s"
	ROUTE_ROOM_CREATE = "/room/create"
)

func Start() {
	router := mux.NewRouter()

	router.HandleFunc(ROUTE_INDEX, index).Methods("GET")
	router.HandleFunc(fmt.Sprintf(ROUTE_ROOM_JOIN, "{room_id}", "{username}"), joinRoom).Methods("GET")
	router.HandleFunc(ROUTE_ROOM_CREATE, createRoom).Methods("GET")

	http.Handle(ROUTE_INDEX, router)

	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "GET"})
	origins := handlers.AllowedOrigins([]string{"*"})

	handlers := handlers.CORS(credentials, methods, origins)(router)
	log.Fatal(http.ListenAndServe(":5000", handlers))
}
