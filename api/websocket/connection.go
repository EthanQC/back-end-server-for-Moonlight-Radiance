// 实时对战同步，玩家状态更新，游戏事件推送
package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Connection 维护单个WebSocket连接
type Connection struct {
	sync.Mutex
	ws       *websocket.Conn
	send     chan []byte
	hub      *Hub
	playerID uint
	roomID   uint
}

func newConnection(ws *websocket.Conn, hub *Hub, playerID, roomID uint) *Connection {
	return &Connection{
		ws:       ws,
		send:     make(chan []byte, 256),
		hub:      hub,
		playerID: playerID,
		roomID:   roomID,
	}
}

// writePump 发送消息到WebSocket连接
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// readPump 从WebSocket连接读取消息
func (c *Connection) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 处理接收到的消息
		c.handleMessage(message)
	}
}

// write 写入消息到WebSocket
func (c *Connection) write(mt int, payload []byte) error {
	c.Lock()
	defer c.Unlock()
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// handleMessage 处理接收到的消息
func (c *Connection) handleMessage(message []byte) {
	var event Event
	if err := json.Unmarshal(message, &event); err != nil {
		return
	}

	// 根据消息类型处理
	switch event.Type {
	case EventCardPlayed:
		c.hub.broadcastToRoom(c.roomID, event)
	case EventTurnChanged:
		c.hub.broadcastToRoom(c.roomID, event)
	}
}
