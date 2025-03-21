package websocket

import (
	"encoding/json"
	"sync"
)

// Hub 管理所有WebSocket连接
type Hub struct {
	sync.RWMutex
	rooms      map[uint]map[*Connection]bool
	register   chan *Connection
	unregister chan *Connection
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uint]map[*Connection]bool),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
	}
}

// Run 运行WebSocket Hub
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.Lock()
			if _, ok := h.rooms[conn.roomID]; !ok {
				h.rooms[conn.roomID] = make(map[*Connection]bool)
			}
			h.rooms[conn.roomID][conn] = true
			h.Unlock()

		case conn := <-h.unregister:
			h.Lock()
			if conns, ok := h.rooms[conn.roomID]; ok {
				if _, ok := conns[conn]; ok {
					delete(conns, conn)
					close(conn.send)
					if len(conns) == 0 {
						delete(h.rooms, conn.roomID)
					}
				}
			}
			h.Unlock()
		}
	}
}

// broadcastToRoom 广播消息到指定房间
func (h *Hub) broadcastToRoom(roomID uint, event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	h.RLock()
	conns := h.rooms[roomID]
	h.RUnlock()

	for conn := range conns {
		select {
		case conn.send <- data:
		default:
			h.unregister <- conn
		}
	}
}
