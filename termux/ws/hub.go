package ws

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Inbound messages from the clients except self.
	broadcastToOthers chan struct {
		From    string
		Message []byte
	}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast: make(chan []byte),
		broadcastToOthers: make(chan struct {
			From    string
			Message []byte
		}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.broadcastToOthers:
			for client := range h.clients {
				if client.conn.RemoteAddr().String() == message.From {
					continue
				}
				select {
				case client.send <- message.Message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Len() int {
	return len(h.clients)
}

func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}
