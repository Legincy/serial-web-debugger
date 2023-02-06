package web

import (
	"github.com/gorilla/websocket"
	"log"
)

type ClientList map[*Client]bool

type Client struct {
	connection   *websocket.Conn
	manager      *Manager
	messageQueue chan []byte
}

func NewClient(connection *websocket.Conn, manager *Manager) *Client {
	return &Client{connection, manager, make(chan []byte)}
}

func (client *Client) readMessages() {
	defer client.manager.removeClient(client)
	socketPackage := &SocketPackage{}

	for {
		err := client.connection.ReadJSON(&socketPackage)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
		log.Println("Received Message:", socketPackage.MessageType, socketPackage.Message)
	}
}

func (client *Client) writeMessage() {
	defer client.manager.removeClient(client)

	for {
		select {
		case message, status := <-client.messageQueue:
			//fmt.Println("Status: ", status, " Message: ", message)
			if !status {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed: ", err)
				}
				return
			}
			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
			}
		}
	}
}
