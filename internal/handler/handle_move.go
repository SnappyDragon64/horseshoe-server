package handler

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"horseshoe-server/internal/packet"
	"horseshoe-server/internal/util"
	"log"
)

func HandleMove(p *game.Player, world *game.World, data []byte) {
	if !p.MoveLimiter.Allow() {
		return
	}

	var req struct {
		Target util.Vector2 `json:"target"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		log.Println("Invalid move packet:", err)
		return
	}

	p.SetPos(req.Target)

	pkt := packet.NewPlayerMovedPacket(p.ID, req.Target)

	if room := p.GetRoom(); room != nil {
		room.Broadcast <- game.BroadcastMsg{
			SenderId: p.ID,
			Data:     pkt,
		}
	}
}
