package main

import (
	"log"
	"net/http"

	"iva/controllers"
	"iva/lib"
	"iva/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	if err := models.EnsureSchema(); err != nil {
		log.Fatal(err)
	}
	// Rotas
	mux := http.NewServeMux()
	mux.HandleFunc("/users/get", controllers.GetUserByEmailHandler)
	mux.HandleFunc("/users/create", controllers.CreateUserHandler)
	mux.HandleFunc("/users/password", controllers.ChangePasswordHandler)
	mux.HandleFunc("/users/edit", controllers.EditUserHandler)
	mux.HandleFunc("/login", controllers.LoginHandler)
	mux.HandleFunc("/chat", controllers.HandleChat)
	// Aplica CORS em todas as rotas
	handler := lib.CORSMiddleware(mux)

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
