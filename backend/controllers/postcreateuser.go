package controllers

import (
	"encoding/json"
	"net/http"
	"iva/models"
)

// Estrutura da requisição
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// Estrutura da resposta
type CreateUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Phone == "" {
		http.Error(w, "Campos obrigatórios ausentes", http.StatusBadRequest)
		return
	}

	// Chama o model
	err := models.CreateUser(req.Name, req.Email, req.Password, req.Phone)
	if err != nil {
		json.NewEncoder(w).Encode(CreateUserResponse{
			Status:  "error",
			Message: "Erro ao criar usuário: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(CreateUserResponse{
		Status:  "success",
		Message: "Usuário criado com sucesso",
	})
}
