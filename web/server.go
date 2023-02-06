package web

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	DataChannel chan []byte
	manager     *Manager
)

type Manager struct {
	clients ClientList
	sync.RWMutex
}

type SocketPackage struct {
	MessageType string `json:"type"`
	Message     string `json:"message"`
}

func NewManager() *Manager {
	return &Manager{clients: make(ClientList)}
}

func Start() {
	manager = NewManager()

	staticServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", staticServer)
	http.HandleFunc("/ws-serial-debugger", manager.serialDebuggerWebsocket)

	log.Println("Starting server at port 8080\n")
	go gatherSerialData()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func gatherSerialData() {
	for {
		if DataChannel == nil {
			continue
		}
		data := <-DataChannel
		for client := range manager.clients {
			client.messageQueue <- data
		}
	}
}

func (manager *Manager) serialDebuggerWebsocket(responseWriter http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := NewClient(connection, manager)
	manager.addClient(client)

	go client.readMessages()
	go client.writeMessage()
}

func (manager *Manager) addClient(client *Client) {
	manager.Lock()
	defer manager.Unlock()

	manager.clients[client] = true
}

func (manager *Manager) removeClient(client *Client) {
	manager.Lock()
	defer manager.Unlock()

	if _, status := manager.clients[client]; status {
		client.connection.Close()
		delete(manager.clients, client)
	}
}
