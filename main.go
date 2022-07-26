package main

import (
	"fmt"

	appHttp "github.com/hypnocill/my-private-chatrooms/api/http"
)

/*
	This project is intended for personal use only.
	In addition to that, the project is currently in an MVP state and needs
	cleanup, architectural changes as well as additional features.
*/
func main() {
	fmt.Println("Private Chatrooms v1 - up and running")
	appHttp.Start()
}
