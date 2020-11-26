package models

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var ws *websocket.Conn

// HandleWebsocket creates the websocket connection
func HandleWebsocket() func(c *fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		ws = c
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	})
}

// Publish writes a message to the websocket
func Publish(message []byte) {
	if ws != nil {
		ws.WriteMessage(1, message)
		log.Printf("wrote %v\n", string(message))
	}
	log.Println("no open websocket connections")
}
