package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: os.Getenv("HOST") + ":8808", Path: "/build"}
	log.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	c.WriteMessage(websocket.TextMessage, []byte("hello"))
	for {
		var bp Blueprint
		err := c.ReadJSON(&bp)
		if err != nil {
			log.Panic("read:", err)
			return
		}
		build(bp)
		fmt.Println(bp.Name + ": done")
		c.WriteMessage(websocket.TextMessage, []byte(bp.Name+": done"))
	}
}
