package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	appHttp "github.com/hypnocill/my-private-chatrooms/api/http"
)

func main() {
	fmt.Println("up and running")
	appHttp.SetupRoutes()

	port := os.Getenv("PORT")
	fmt.Println("port")
	fmt.Println(port)
	if port == "" {
		port = "5000"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
