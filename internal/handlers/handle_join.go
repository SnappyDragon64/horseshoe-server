package handlers

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"horseshoe-server/internal/utils"
	"log"
)

func HandleJoin(p *game.Player, world *game.World, data []byte) {
	var req struct {
		Room string        `json:"room"`
		Pos  utils.Vector2 `json:"pos"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		log.Println("Invalid join packet:", err)
		return
	}

	targetRoom, exists := world.Rooms[req.Room]
	if !exists {
		return
	}

	if oldRoom := p.GetRoom(); oldRoom != nil {
		oldRoom.Leave <- p
	}

	p.SetPos(req.Pos)

	targetRoom.Join <- p
}
