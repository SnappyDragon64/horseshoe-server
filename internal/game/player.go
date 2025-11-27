package game

import (
	"horseshoe-server/internal/utils"
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

type Player struct {
	ID   string
	Conn *websocket.Conn
	room *Room
	pos  utils.Vector2
	mu   sync.RWMutex

	Send chan []byte
}

func NewPlayer(id string, conn *websocket.Conn) *Player {
	return &Player{
		ID:   id,
		Conn: conn,
		Send: make(chan []byte, 1024),
	}
}

func (p *Player) SetRoom(r *Room) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.room = r
}

func (p *Player) GetRoom() *Room {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.room
}

func (p *Player) SetPos(pos utils.Vector2) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pos = pos
}

func (p *Player) GetPos() utils.Vector2 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.pos
}

func (p *Player) ReadPump(handleMessage func([]byte)) {
	defer func() {
		p.Conn.Close()
	}()

	p.Conn.SetReadLimit(maxMessageSize)

	p.Conn.SetReadDeadline(time.Now().Add(pongWait))

	p.Conn.SetPongHandler(func(string) error {
		p.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		handleMessage(message)
	}
}

func (p *Player) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-p.Send:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(p.Send)
			for i := 0; i < n; i++ {
				w.Write(<-p.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
