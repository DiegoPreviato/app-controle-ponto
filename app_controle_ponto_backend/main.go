package main

import (
	"fmt"
	"log"
	"net/http"

	"controle-ponto-api/database"
	"controle-ponto-api/handlers"
	"controle-ponto-api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	fmt.Println("Starting Ponto Control API...")

	// Load environment variables (e.g., from a .env file) would be a good addition here
	// For now, it relies on system-set env vars

	if err := database.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.DB.Close()

	r := chi.NewRouter()

	// CORS Middleware
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all for dev
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is live!"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public auth routes
		r.Post("/register", handlers.Register)
		r.Post("/login", handlers.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.JwtAuthentication)

			r.Post("/pontos", handlers.RegistrarPonto)
			r.Get("/pontos/{data}", handlers.ListarPontosPorData)
			r.Get("/pontos/{data}/total-horas", handlers.CalcularHorasTrabalhadas)
			r.Put("/pontos/{id}", handlers.AtualizarPonto)
			r.Delete("/pontos/{id}", handlers.DeletarPonto)
		})
	})

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
