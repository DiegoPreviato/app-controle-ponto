package main

import (
	"fmt"
	"log"
	"net/http"

	"controle-ponto-api/database"
	"controle-ponto-api/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors" // Import the cors package
)

func main() {
	fmt.Println("Iniciando a API de Controle de Ponto...")

	if err := database.InitDB("./ponto.db"); err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}
	defer database.DB.Close()

	r := chi.NewRouter()

	// Configure CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for development
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(corsMiddleware.Handler) // Use the CORS middleware

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API no ar!"))
	})

	r.Post("/registrar-ponto", handlers.RegistrarPonto)
	r.Get("/pontos/{data}", handlers.ListarPontosPorData)
	r.Get("/pontos/{data}/total-horas", handlers.CalcularHorasTrabalhadas)

	fmt.Println("Servidor escutando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}