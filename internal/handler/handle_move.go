package handler

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"horseshoe-server/internal/packet"
	"horseshoe-server/internal/util"
	"math"
)

func HandleMove(p *game.Player, world *game.World, data []byte) {
	if !p.MoveLimiter.Allow() {
		return
	}

	var req struct {
		Target util.Vector2 `json:"target"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return
	}

	if math.Abs(req.Target.X) > 10000 || math.Abs(req.Target.Y) > 10000 {
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
