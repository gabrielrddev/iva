package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Estruturas de requisição para a API Gemini (Google AI)
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Role  string `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

// Estruturas de resposta para a API Gemini
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []GeminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Função que envia a pergunta para o Gemini
func AskGemini(question string) (string, error) {
	// 1. Configuração e Chave de API
	apiKey := os.Getenv("GEMINI_API_KEY") 
	if apiKey == "" {
		return "", fmt.Errorf("variável de ambiente GEMINI_API_KEY não configurada")
	}

	// Mantendo a variável de ambiente SOURCES original
	sources := os.Getenv("GPT_SOURCES")
	
	// 2. Instrução e prompt unificados no conteúdo do usuário (sem campo config)
	systemInstruction := "Você é um assistente técnico e comercial especializado em produtos da Intelbras chamado I.V.A (Intelbras Virtual Assistant). Todas as suas respostas devem se basear EXCLUSIVAMENTE em produtos, manuais, e informações oficiais da Intelbras. Não forneça informações sobre produtos de outras marcas ou assuntos fora desse contexto. Sempre responda com base em produtos Intelbras, e jamais mencione marcas concorrentes e você não precisa se apresentar o usuario ja sabe quem você é."
	userPrompt := fmt.Sprintf("%s\n\nUse como referência os seguintes links e fontes oficiais (fornecidos pelo sistema):\n%s\n\nPergunta do usuário:\n%s\n", systemInstruction, sources, question)

	// 3. Montagem do Corpo da Requisição (apenas contents)
	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Role: "user",
				Parts: []GeminiPart{{Text: userPrompt}},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar JSON: %w", err)
	}
	
	// 4. Endpoint e Requisição HTTP
	model := "gemini-2.5-flash"
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// 5. Execução da Requisição
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// 6. Tratamento de Erros HTTP
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("Gemini API HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(b)))
	}

	// 7. Decodificação da Resposta
	var geminiRes GeminiResponse
	if err := json.NewDecoder(res.Body).Decode(&geminiRes); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta da Gemini: %w", err)
	}

	// 8. Extração do Conteúdo
	if len(geminiRes.Candidates) > 0 && len(geminiRes.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(geminiRes.Candidates[0].Content.Parts[0].Text), nil
	}

	return "Não foi possível obter resposta da IA Gemini. Verifique se o conteúdo da pergunta é apropriado.", nil
}

// Wrapper para compatibilidade com chamadas existentes
func AskChatGPT(question string) (string, error) {
	return AskGemini(question)
} 