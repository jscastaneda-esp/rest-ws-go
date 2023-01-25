package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	id       string
	socket   *websocket.Conn
	outbound chan []byte
}

func (client *Client) Write() {
	for {
		select {
		case message, ok := <-client.outbound:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			client.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		Hub:      hub,
		id:       socket.RemoteAddr().String(),
		socket:   socket,
		outbound: make(chan []byte),
	}
}
