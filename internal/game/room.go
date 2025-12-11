package game

import (
	"encoding/json"
	"horseshoe-server/internal/packet"
	"log"
)

type BroadcastMsg struct {
	SenderId string
	Data     interface{}
}

type Room struct {
	ID      string
	Players map[string]*Player

	Join      chan *Player
	Leave     chan *Player
	Broadcast chan BroadcastMsg
}

func NewRoom(id string) *Room {
	return &Room{
		ID:        id,
		Players:   make(map[string]*Player),
		Join:      make(chan *Player),
		Leave:     make(chan *Player),
		Broadcast: make(chan BroadcastMsg),
	}
}

func (r *Room) Run() {
	for {
		select {
		case p := <-r.Join:
			r.Players[p.ID] = p
			p.SetRoom(r)

			currentPlayersList := make([]packet.PlayerData, 0)

			for _, existingP := range r.Players {
				if existingP.ID == p.ID {
					continue
				}

				currentPlayersList = append(currentPlayersList, packet.PlayerData{
					ID:  existingP.ID,
					Pos: existingP.GetPos(),
				})

				spawnPkt := packet.NewSpawnPlayerPacket(p.ID, p.GetPos())
				existingP.SendPacket(spawnPkt)
			}

			loadPkt := packet.NewLoadRoomPacket(r.ID, p.GetPos(), currentPlayersList)
			p.SendPacket(loadPkt)

		case p := <-r.Leave:
			if _, exists := r.Players[p.ID]; exists {
				delete(r.Players, p.ID)
				if p.GetRoom() == r {
					p.SetRoom(nil)
				}

				leavePkt := packet.NewDeletePlayerPacket(p.ID)
				for _, p_ := range r.Players {
					p_.SendPacket(leavePkt)
				}
			}

		case msg := <-r.Broadcast:
			bytes, err := json.Marshal(msg.Data)
			if err != nil {
				log.Println("Broadcast marshal error:", err)
				continue
			}

			for id, p := range r.Players {
				if id != msg.SenderId {
					select {
					case p.Send <- bytes:
					default:
						p.Disconnect()
					}
				}
			}
		}
	}
}
