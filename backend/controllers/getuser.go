package controllers

import (
	"iva/models"
	"encoding/json"
	"net/http"
)

// Estrutura de resposta do usuário (sem senha)
type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// GET /users?email=exemplo@email.com
func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Recebe o email como query param
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email do usuário é obrigatório", http.StatusBadRequest)
		return
	}

	// Busca usuário no model
	user, err := models.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	resp := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}

	jsonResponse(w, resp)
}

// Função helper para JSON (reutilizável)
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
