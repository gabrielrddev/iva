package controllers

import (
	"iva/models"
	"encoding/json"
	"net/http"
)

// Estrutura da requisição de edição do usuário
type EditUserRequest struct {
	Email    string `json:"email"`    // Email atual do usuário (para identificar)
	NewName  string `json:"newName"`  // Novo nome
	NewEmail string `json:"newEmail"` // Novo email
	NewPhone string `json:"newPhone"` // Novo telefone
}

// Estrutura de resposta
type EditUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// POST /users/edit
func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req EditUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Busca usuário pelo email atual
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		jsonResponse(w, EditUserResponse{
			Status:  "error",
			Message: "Usuário não encontrado",
		})
		return
	}

	// Atualiza usuário no banco
	if err := models.EditUser(user.Email, req.NewName, req.NewEmail, req.NewPhone); err != nil {
		jsonResponse(w, EditUserResponse{
			Status:  "error",
			Message: "Erro ao atualizar usuário",
		})
		return
	}

	jsonResponse(w, EditUserResponse{Status: "success"})
}
