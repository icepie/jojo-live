package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	WsHub *Hub
)

func init() {
	WsHub = NewHub()
}

type WsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (m WsMessage) ToJson() []byte {
	d, _ := json.Marshal(m)
	return d
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
