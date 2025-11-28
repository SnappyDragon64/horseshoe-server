package packets

import (
	"horseshoe-server/internal/utils"
)

type ConnectPacket struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func NewConnectPacket(id string) *ConnectPacket {
	return &ConnectPacket{Type: "connect", ID: id}
}

type PlayerData struct {
	ID  string        `json:"id"`
	Pos utils.Vector2 `json:"pos"`
}

type LoadRoomPacket struct {
	Type    string        `json:"type"`
	RoomID  string        `json:"id"`
	Pos     utils.Vector2 `json:"pos"`
	Players []PlayerData  `json:"players"`
}

func NewLoadRoomPacket(roomID string, pos utils.Vector2, players []PlayerData) *LoadRoomPacket {
	return &LoadRoomPacket{
		Type:    "load_room",
		RoomID:  roomID,
		Pos:     pos,
		Players: players,
	}
}

type SpawnPlayerPacket struct {
	Type string        `json:"type"`
	ID   string        `json:"id"`
	Pos  utils.Vector2 `json:"pos"`
}

func NewSpawnPlayerPacket(id string, pos utils.Vector2) *SpawnPlayerPacket {
	return &SpawnPlayerPacket{Type: "spawn_player", ID: id, Pos: pos}
}

type DeletePlayerPacket struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func NewDeletePlayerPacket(id string) *DeletePlayerPacket {
	return &DeletePlayerPacket{Type: "delete_player", ID: id}
}

type PlayerMovedPacket struct {
	Type   string        `json:"type"`
	ID     string        `json:"id"`
	Target utils.Vector2 `json:"target"`
}

func NewPlayerMovedPacket(id string, target utils.Vector2) *PlayerMovedPacket {
	return &PlayerMovedPacket{Type: "player_moved", ID: id, Target: target}
}

type PlayerMessagePacket struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

func NewPlayerMessagePacket(id, msg string) *PlayerMessagePacket {
	return &PlayerMessagePacket{Type: "player_message", ID: id, Message: msg}
}
