package handlers

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"log"
)

func HandleChat(p *game.Player, world *game.World, data []byte) {
	var req struct {
		Message string `json:"message"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		log.Println("Invalid chat packet:", err)
		return
	}

	msg := []rune(req.Message)
	if len(msg) > 120 {
		msg = msg[:120]
	}

	if len(msg) == 0 {
		return
	}

	broadcastData, _ := json.Marshal(map[string]interface{}{
		"type":    "player_message",
		"id":      p.ID,
		"message": string(msg),
	})

	if room := p.GetRoom(); room != nil {
		room.Broadcast <- game.BroadcastMsg{
			SenderId: p.ID,
			Data:     broadcastData,
		}
	}
}
