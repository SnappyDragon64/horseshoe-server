package handlers

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"log"
)

type HandlerFunc func(p *game.Player, world *game.World, data []byte)

var PacketHandlers = map[string]HandlerFunc{
	"join_room": HandleJoin,
	"move":      HandleMove,
	"chat":      HandleChat,
}

func ProcessPacket(p *game.Player, world *game.World, data []byte) {
	var header struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &header); err != nil {
		log.Println("Invalid JSON header:", err)
		return
	}

	if handler, exists := PacketHandlers[header.Type]; exists {
		handler(p, world, data)
	} else {
		log.Printf("Unknown packet type: %s", header.Type)
	}
}
