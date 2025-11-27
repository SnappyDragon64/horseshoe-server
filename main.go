package main

import (
	"encoding/json"
	"horseshoe-server/internal/game"
	"horseshoe-server/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	world := game.NewWorld()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade failed:", err)
			return
		}

		newId := world.GeneratePlayerId()
		player := game.NewPlayer(newId, conn)
		world.AddPlayer(player)

		connectMsg, _ := json.Marshal(map[string]string{
			"type": "connect",
			"id":   newId,
		})
		player.Send <- connectMsg

		go player.WritePump()

		player.ReadPump(func(message []byte) {
			handlers.ProcessPacket(player, world, message)
		})

		if room := player.GetRoom(); room != nil {
			room.Leave <- player
		}
		world.RemovePlayer(newId)
	})

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server start error: ", err)
	}
}
