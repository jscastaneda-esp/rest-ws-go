package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("Client Connected", client.socket.RemoteAddr())

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client Disconnected", client.socket.RemoteAddr())
	client.socket.Close()

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	indexDelete := -1
	for index, cli := range hub.clients {
		if cli.id == client.id {
			indexDelete = index
			break
		}
	}

	hub.clients = append(hub.clients[:indexDelete], hub.clients[indexDelete+1:]...)
	// copy(hub.clients[indexDelete:], hub.clients[indexDelete+1:])
	// hub.clients[len(hub.clients)-1] = nil
	// hub.clients = hub.clients[:len(hub.clients)-1]
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	client := NewClient(hub, socket)
	hub.register <- client

	go client.Write()
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) Broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ignore {
			client.outbound <- data
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}
