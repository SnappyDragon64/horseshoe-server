package api

import (
	"encoding/json"
	"horseshoe-server/internal/auth"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
}
