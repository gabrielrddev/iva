package controllers

import (
	"iva/models"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Estrutura da requisição de mudança de senha
type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Estrutura de resposta
type ChangePasswordResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// POST /users/password
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Busca usuário pelo email
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		jsonResponse(w, ChangePasswordResponse{
			Status:  "error",
			Message: "Usuário não encontrado",
		})
		return
	}

	// Verifica senha antiga
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		jsonResponse(w, ChangePasswordResponse{
			Status:  "error",
			Message: "Senha antiga incorreta",
		})
		return
	}

	// Atualiza senha
	if err := models.UpdateUserPassword(user.Email, req.NewPassword); err != nil {
		jsonResponse(w, ChangePasswordResponse{
			Status:  "error",
			Message: "Erro ao atualizar senha",
		})
		return
	}

	jsonResponse(w, ChangePasswordResponse{Status: "success"})
}
