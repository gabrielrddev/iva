package models

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Phone    string
	Questions int
}

func IncrementUserQuestions(email string) error {
	db, err := sql.Open("sqlite3", os.Getenv("PATH_DB"))
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET questions = questions + 1 WHERE email = ?", email)
	return err
}

func CreateUser(name, email, password, phone string) error {
	db, err := sql.Open("sqlite3", os.Getenv("PATH_DB"))
	if err != nil {
		log.Println("Erro ao abrir banco:", err)
		return err
	}
	defer db.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (name, email, password, phone) VALUES (?, ?, ?, ?)",
		name, email, string(hashedPassword), phone)
	if err != nil {
		return err
	}

	return nil
}

// Busca usuário por email
func GetUserByEmail(email string) (*User, error) {
	db, err := sql.Open("sqlite3", os.Getenv("PATH_DB"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	var u User
	err = db.QueryRow("SELECT id, name, phone, email, password FROM users WHERE email = ?", email).
		Scan(&u.ID, &u.Name, &u.Phone, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func UpdateUserPassword(email, newPassword string) error {
	db, err := sql.Open("sqlite3", os.Getenv("PATH_DB"))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	// Hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Atualiza no banco
	_, err = db.Exec("UPDATE users SET password = ? WHERE email = ?", string(hashedPassword), email)
	if err != nil {
		return err
	}

	return nil
}

func EditUser(currentEmail, newName, newEmail, newPhone string) error {
	db, err := sql.Open("sqlite3", os.Getenv("PATH_DB"))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		"UPDATE users SET name = ?, email = ?, phone = ? WHERE email = ?",
		newName, newEmail, newPhone, currentEmail,
	)
	if err != nil {
		return err
	}

	return nil
}

func EnsureSchema() error {
	db, err := sql.Open("sqlite3", os.Getenv("PATH_DB"))
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        phone TEXT NOT NULL,
        questions INTEGER DEFAULT 0
    )`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`ALTER TABLE users ADD COLUMN questions INTEGER DEFAULT 0`)
	if err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
			return err
		}
	}

	return nil
}
