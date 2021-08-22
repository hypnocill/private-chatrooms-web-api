package main

import (
	"fmt"

	appHttp "github.com/hypnocill/my-private-chatrooms/api/http"
)

func main() {
	fmt.Println("up and running")
	appHttp.Start()
}
