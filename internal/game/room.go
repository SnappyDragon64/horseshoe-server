package game

import (
	"encoding/json"
)

type BroadcastMsg struct {
	SenderId string
	Data     []byte
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

			currentPlayersList := make([]map[string]interface{}, 0)

			for _, existingP := range r.Players {
				if existingP.ID == p.ID {
					continue
				}

				currentPlayersList = append(currentPlayersList, map[string]interface{}{
					"id":  existingP.ID,
					"pos": existingP.GetPos(),
				})

				spawnPkt, _ := json.Marshal(map[string]interface{}{
					"type": "spawn_player",
					"id":   p.ID,
					"pos":  p.GetPos(),
				})

				select {
				case existingP.Send <- spawnPkt:
				default:
					close(existingP.Send)
					delete(r.Players, existingP.ID)
				}
			}

			loadPkt, _ := json.Marshal(map[string]interface{}{
				"type":    "load_room",
				"id":      r.ID,
				"pos":     p.GetPos(),
				"players": currentPlayersList,
			})
			p.Send <- loadPkt

		case p := <-r.Leave:
			if _, exists := r.Players[p.ID]; exists {
				delete(r.Players, p.ID)
				if p.GetRoom() == r {
					p.SetRoom(nil)
				}

				leavePkt, _ := json.Marshal(map[string]string{
					"type": "delete_player",
					"id":   p.ID,
				})

				for _, p_ := range r.Players {
					select {
					case p_.Send <- leavePkt:
					default:
						close(p_.Send)
						delete(r.Players, p_.ID)
					}
				}
			}

		case msg := <-r.Broadcast:
			for id, p := range r.Players {
				if id != msg.SenderId {
					select {
					case p.Send <- msg.Data:
					default:
						close(p.Send)
						delete(r.Players, id)
					}
				}
			}
		}
	}
}
