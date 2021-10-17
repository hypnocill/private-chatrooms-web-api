package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	ROUTE_INDEX       = "/"
	ROUTE_ROOM_AUTH   = "/room/auth"
	ROUTE_ROOM_JOIN   = "/room/join/%s/%s"
	ROUTE_ROOM_CREATE = "/room/create"
)

func Start() {
	router := mux.NewRouter()

	router.HandleFunc(ROUTE_INDEX, index).Methods("GET")
	router.HandleFunc(ROUTE_ROOM_AUTH, authRoom).Methods("POST")
	router.HandleFunc(fmt.Sprintf(ROUTE_ROOM_JOIN, "{room_id}", "{username}"), joinRoom).Methods("GET")
	router.HandleFunc(ROUTE_ROOM_CREATE, createRoom).Methods("POST")

	http.Handle(ROUTE_INDEX, router)

	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "GET"})
	origins := handlers.AllowedOrigins([]string{"*"})

	handlers := handlers.CORS(credentials, methods, origins)(router)

	port := ":" + getPort()
	log.Fatal(http.ListenAndServe(port, handlers))
}

func getPort() string {
	port := "5000"
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	return port
}
