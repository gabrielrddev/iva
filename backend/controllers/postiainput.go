package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"log"
	"strings"
	"iva/services"

	_ "github.com/mattn/go-sqlite3"
)

type ChatRequest struct {
	Email   string `json:"email"`
	Input   string `json:"input"`
	Message string `json:"message"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// Controller que processa a pergunta e envia para o GPT
func HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}

	if req.Input == "" && req.Message != "" {
		req.Input = req.Message
	}
	if req.Input == "" {
		http.Error(w, "Campo 'input' é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) != "" {
		db, err := sql.Open("sqlite3", os.Getenv("PATH_DB"))
		if err != nil {
			http.Error(w, "Erro no banco", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Incrementa o número de perguntas do usuário
		_, err = db.Exec("UPDATE users SET questions = questions + 1 WHERE email = ?", req.Email)
		if err != nil {
			http.Error(w, "Erro ao atualizar número de perguntas", http.StatusInternalServerError)
			return
		}
	}

	// Chama o service para enviar a pergunta ao ChatGPT
	answer, err := services.AskChatGPT(req.Input)
	if err != nil {
		log.Printf("Erro ao consultar ChatGPT: %v", err)
		http.Error(w, "Erro ao consultar ChatGPT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Response: answer})
}
