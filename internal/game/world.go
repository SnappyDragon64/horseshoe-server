package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
)

type World struct {
	Players map[string]*Player
	Rooms   map[string]*Room
	mu      sync.Mutex
}

func NewWorld() *World {
	w := &World{
		Players: make(map[string]*Player),
		Rooms:   make(map[string]*Room),
	}

	file, err := os.ReadFile("config/rooms.json")
	if err != nil {
		panic(err)
	}

	var roomNames []string
	err = json.Unmarshal(file, &roomNames)
	if err != nil {
		panic(err)
	}

	for _, name := range roomNames {
		r := NewRoom(name)
		w.Rooms[name] = r

		go r.Run()
	}

	return w
}

func (w *World) GeneratePlayerId() string {
	w.mu.Lock()
	defer w.mu.Unlock()

	for {
		id := fmt.Sprintf("P%d", 10000000+rand.Intn(90000000))

		if _, exists := w.Players[id]; !exists {
			return id
		}
	}
}

func (w *World) AddPlayer(p *Player) {
	w.mu.Lock()
	w.Players[p.ID] = p
	w.mu.Unlock()
}

func (w *World) RemovePlayer(id string) {
	w.mu.Lock()
	delete(w.Players, id)
	w.mu.Unlock()
}
