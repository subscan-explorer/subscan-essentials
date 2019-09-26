package websocket

type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: //join
			h.clients[client] = true
			go client.join()
		case client := <-h.unregister: //exit
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				go client.leave()
				close(client.send)
			}
		case message := <-h.Broadcast: //broadcast
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
