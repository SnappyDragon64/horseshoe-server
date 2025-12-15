package handler

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"log"
)

func HandleJoin(p *game.Player, world *game.World, data []byte) {
	var req struct {
		Room string `json:"room"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		log.Println("Invalid join packet:", err)
		return
	}

	targetRoom, exists := world.Rooms[req.Room]
	if !exists {
		log.Printf("Player %s attempted to join invalid room: %s", p.ID, req.Room)
		return
	}

	currentRoom := p.GetRoom()

	originKey := "default"
	if currentRoom != nil {
		originKey = currentRoom.ID
	}

	roomConfig := world.RoomsConfig[req.Room]

	spawnPos, found := roomConfig.Spawns[originKey]
	if !found {
		spawnPos = roomConfig.Spawns["default"]
	}

	if currentRoom != nil {
		currentRoom.Leave <- p
	}

	p.SetPos(spawnPos)

	targetRoom.Join <- p
}
