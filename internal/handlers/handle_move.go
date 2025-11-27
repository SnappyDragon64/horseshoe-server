package handlers

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"horseshoe-server/internal/utils"
	"log"
)

func HandleMove(p *game.Player, world *game.World, data []byte) {
	var req struct {
		Target utils.Vector2 `json:"target"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		log.Println("Invalid move packet:", err)
		return
	}

	p.SetPos(req.Target)

	broadcastData, _ := json.Marshal(map[string]interface{}{
		"type":   "player_moved",
		"id":     p.ID,
		"target": req.Target,
	})

	if room := p.GetRoom(); room != nil {
		room.Broadcast <- game.BroadcastMsg{
			SenderId: p.ID,
			Data:     broadcastData,
		}
	}
}
