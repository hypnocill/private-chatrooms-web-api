package main

import (
	"fmt"
	"log"
	"net/http"

	appHttp "github.com/hypnocill/my-private-chatrooms/api/http"
)

func main() {
	fmt.Println("up and running")
	appHttp.SetupRoutes()
	log.Fatal(http.ListenAndServe(":5000", nil))
}
