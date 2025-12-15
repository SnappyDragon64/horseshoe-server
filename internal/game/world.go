package game

import (
	"encoding/json"
	"horseshoe-server/internal/util"
	"log"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type RoomConfig struct {
	Spawns map[string]util.Vector2 `json:"spawns"`
}

type World struct {
	Players     map[string]*Player
	Rooms       map[string]*Room
	RoomsConfig map[string]RoomConfig
	mu          sync.Mutex
}

func NewWorld() *World {
	w := &World{
		Players:     make(map[string]*Player),
		Rooms:       make(map[string]*Room),
		RoomsConfig: make(map[string]RoomConfig),
	}

	file, err := os.ReadFile("config/rooms.json")
	if err != nil {
		log.Fatalf("Failed to load rooms config: %v", err)
	}

	err = json.Unmarshal(file, &w.RoomsConfig)
	if err != nil {
		log.Fatalf("Failed to parse rooms config: %v", err)
	}

	for id := range w.RoomsConfig {
		r := NewRoom(id)
		w.Rooms[id] = r

		go r.Run()
	}

	return w
}

func (w *World) AddPlayer(p *Player) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if oldPlayer, exists := w.Players[p.ID]; exists {
		msg := websocket.FormatCloseMessage(4000, "Logged in from another location")
		oldPlayer.Conn.WriteMessage(websocket.CloseMessage, msg)
		oldPlayer.Conn.Close()
	}

	w.Players[p.ID] = p
}

func (w *World) RemovePlayer(p *Player) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if currentP, exists := w.Players[p.ID]; exists {
		if currentP == p {
			delete(w.Players, p.ID)
		}
	}
}
