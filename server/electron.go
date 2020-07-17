package server

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Electron struct {
	hub  *Hub
	send chan []byte
	conn *websocket.Conn
	id   string
}

func (e *Electron) readMessages() {
	defer func() {
		e.hub.disconnect <- e
		e.conn.Close()
	}()

	for {
		_, message, err := e.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		fmt.Println("Message Recieved", string(message))

		e.send <- []byte("Recieved Your Message")

	}
}

func (e *Electron) writeMessagesToClient() {
	for {
		select {
		case msg := <-e.send:
			e.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
