package handler

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"horseshoe-server/internal/packet"
	"log"
)

func HandleChat(p *game.Player, world *game.World, data []byte) {
	if !p.ChatLimiter.Allow() {
		p.SendPacket(packet.NewSystemMessagePacket("Slow down! You are chatting too fast."))
		return
	}

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

	pkt := packet.NewPlayerMessagePacket(p.ID, string(msg))

	if room := p.GetRoom(); room != nil {
		room.Broadcast <- game.BroadcastMsg{
			SenderId: p.ID,
			Data:     pkt,
		}
	}
}
