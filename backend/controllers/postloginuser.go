package controllers

import (
	"iva/models"
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

// Estrutura para receber login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Estrutura de resposta
type LoginResponse struct {
	Status string `json:"status"`
}

// /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Busca usuário usando função do model
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		jsonResponse(w, LoginResponse{Status: "user_not_found"})
		return
	}

	// Verifica senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		jsonResponse(w, LoginResponse{Status: "wrong_password"})
		return
	}

	jsonResponse(w, LoginResponse{Status: "success"})
}
