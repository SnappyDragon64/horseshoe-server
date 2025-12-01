package main

import (
	"encoding/json"
	"fmt"
	"horseshoe-server/internal/auth"
	"horseshoe-server/internal/db"
	"horseshoe-server/internal/game"
	"horseshoe-server/internal/handlers"
	"horseshoe-server/internal/packets"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"https://html.itch.zone": true,
			"http://localhost:8060":  true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "local_dev.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Println("WARNING: JWT_SECRET not set. Using default 'dev-secret'.")
		jwtSecret = "dev-secret"
	}

	auth.JwtSecret = []byte(jwtSecret)

	db.Init(dbPath)

	world := game.NewWorld()

	http.HandleFunc("/api/register", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad JSON", 400)
			return
		}

		if err := auth.Register(req.Username, req.Password); err != nil {
			http.Error(w, "Registration failed: "+err.Error(), 409)
			return
		}

		w.Write([]byte(`{"status":"ok"}`))
	}))

	http.HandleFunc("/api/login", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		token, username, err := auth.Login(req.Username, req.Password)
		if err != nil {
			http.Error(w, "Invalid credentials", 401)
			return
		}

		response := map[string]string{
			"token":    token,
			"username": username,
		}
		json.NewEncoder(w).Encode(response)
	}))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			http.Error(w, "Missing authentication token", 401)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return auth.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", 401)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)

		var user db.User
		result := db.DB.Select("id").Where("username = ?", username).First(&user)

		if result.Error != nil {
			http.Error(w, "User no longer exists or DB error", 401)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade failed:", err)
			return
		}

		player := game.NewPlayer(username, conn)
		world.AddPlayer(player)

		connectPkt := packets.NewConnectPacket(username)
		player.SendPacket(connectPkt)

		go player.WritePump()

		player.ReadPump(func(message []byte) {
			handlers.ProcessPacket(player, world, message)
		})

		if room := player.GetRoom(); room != nil {
			room.Leave <- player
		}

		world.RemovePlayer(player)
	})

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server start error: ", err)
	}
}
