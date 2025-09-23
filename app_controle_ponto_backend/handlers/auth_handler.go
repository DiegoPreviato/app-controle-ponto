package handlers

import (
	"controle-ponto-api/database"
	"controle-ponto-api/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// Register godoc
// @Summary      Registra um novo usuário
// @Description  Cria um novo usuário no sistema com nome, email e senha.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "Dados do usuário para registro (ID e Horarios podem ser omitidos)"
// @Success      201   {object}  map[string]string
// @Failure      400   {string}  string "Invalid request body"
// @Failure      500   {string}  string "Failed to create user"
// @Router       /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (nome, email, password_hash) VALUES ($1, $2, $3)", user.Nome, user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// Login godoc
// @Summary      Realiza o login do usuário
// @Description  Autentica um usuário com email e senha e retorna um token JWT.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.User  true  "Credenciais de login (apenas email e password são necessários)"
// @Success      200          {object}  map[string]string
// @Failure      400          {string}  string "Invalid request body"
// @Failure      401          {string}  string "Invalid credentials"
// @Failure      500          {string}  string "Internal server error"
// @Router       /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, nome, email, password_hash FROM users WHERE email = $1", creds.Email).Scan(&user.ID, &user.Nome, &user.Email, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}