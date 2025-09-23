package handlers

import (
	"controle-ponto-api/database"
	"controle-ponto-api/middleware"
	"controle-ponto-api/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// --- Funções Auxiliares ---

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// PontoUpdatePayload define a estrutura para o corpo da requisição de atualização de ponto.
// Isso é usado apenas para a documentação do Swagger.
type PontoUpdatePayload struct {
	Horario time.Time `json:"horario"`
}

// --- Handlers ---

// RegistrarPonto godoc
// @Summary      Registra um novo ponto
// @Description  Cria um novo registro de ponto com o horário atual para o usuário autenticado.
// @Tags         Pontos
// @Produce      json
// @Security     ApiKeyAuth
// @Success      201  {object}  models.Ponto
// @Failure      500  {string}  string "Internal server error"
// @Router       /pontos [post]
func RegistrarPonto(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user ID from context")
		return
	}

	horarioDoPonto := time.Now()
	novoPonto := models.Ponto{
		UserID:  userID,
		Horario: horarioDoPonto,
	}

	err := database.DB.QueryRow(
		"INSERT INTO pontos(user_id, horario) VALUES($1, $2) RETURNING id",
		userID, horarioDoPonto,
	).Scan(&novoPonto.ID)

	if err != nil {
		log.Printf("Error inserting new 'ponto': %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to register 'ponto'")
		return
	}

	respondWithJSON(w, http.StatusCreated, novoPonto)
}

// ListarPontosPorData godoc
// @Summary      Lista os pontos por data
// @Description  Lista todos os registros de ponto de um usuário para uma data específica.
// @Tags         Pontos
// @Produce      json
// @Security     ApiKeyAuth
// @Param        data  path      string  true  "Data no formato YYYY-MM-DD"
// @Success      200   {array}   models.Ponto
// @Failure      400   {string}  string  "Invalid date format. Use YYYY-MM-DD"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /pontos/{data} [get]
func ListarPontosPorData(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user ID from context")
		return
	}

	dataParam := chi.URLParam(r, "data")
	parsedDate, err := time.Parse("2006-01-02", dataParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	startOfDay := parsedDate
	endOfDay := startOfDay.Add(24 * time.Hour)

	rows, err := database.DB.Query(
		"SELECT id, user_id, horario FROM pontos WHERE user_id = $1 AND horario >= $2 AND horario < $3 ORDER BY horario ASC",
		userID, startOfDay, endOfDay,
	)
	if err != nil {
		log.Printf("Error querying 'pontos' by date: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve 'pontos'")
		return
	}
	defer rows.Close()

	pontos := []models.Ponto{}
	for rows.Next() {
		var p models.Ponto
		if err := rows.Scan(&p.ID, &p.UserID, &p.Horario); err != nil {
			log.Printf("Error scanning 'ponto' row: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to process 'pontos' data")
			return
		}
		pontos = append(pontos, p)
	}

	respondWithJSON(w, http.StatusOK, pontos)
}

// CalcularHorasTrabalhadas godoc
// @Summary      Calcula horas trabalhadas
// @Description  Calcula o total de horas trabalhadas em um dia com base nos registros de ponto (entrada/saída).
// @Tags         Pontos
// @Produce      json
// @Security     ApiKeyAuth
// @Param        data  path      string  true  "Data no formato YYYY-MM-DD"
// @Success      200   {object}  map[string]string
// @Failure      400   {string}  string  "Invalid date format. Use YYYY-MM-DD"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /pontos/{data}/total-horas [get]
func CalcularHorasTrabalhadas(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user ID from context")
		return
	}

	dataParam := chi.URLParam(r, "data")
	parsedDate, err := time.Parse("2006-01-02", dataParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	startOfDay := parsedDate
	endOfDay := startOfDay.Add(24 * time.Hour)

	rows, err := database.DB.Query(
		"SELECT horario FROM pontos WHERE user_id = $1 AND horario >= $2 AND horario < $3 ORDER BY horario ASC",
		userID, startOfDay, endOfDay,
	)
	if err != nil {
		log.Printf("Error querying 'pontos' for calculation: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve 'pontos' for calculation")
		return
	}
	defer rows.Close()

	var horarios []time.Time
	for rows.Next() {
		var horario time.Time
		if err := rows.Scan(&horario); err != nil {
			log.Printf("Error scanning 'horario' for calculation: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to process 'horarios' for calculation")
			return
		}
		horarios = append(horarios, horario)
	}

	var totalDuracao time.Duration
	if len(horarios)%2 != 0 {
		horarios = horarios[:len(horarios)-1]
	}

	for i := 0; i < len(horarios); i += 2 {
		entrada := horarios[i]
		saida := horarios[i+1]
		totalDuracao += saida.Sub(entrada)
	}

	totalHoras := int(totalDuracao.Hours())
	totalMinutos := int(totalDuracao.Minutes()) % 60

	resposta := map[string]string{
		"total_trabalhado": fmt.Sprintf("%dh %dm", totalHoras, totalMinutos),
		"total_segundos":   fmt.Sprintf("%.0f", totalDuracao.Seconds()),
	}

	respondWithJSON(w, http.StatusOK, resposta)
}

// AtualizarPonto godoc
// @Summary      Atualiza um registro de ponto
// @Description  Atualiza o horário de um registro de ponto existente.
// @Tags         Pontos
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      int                  true  "ID do Ponto"
// @Param        horario  body      PontoUpdatePayload   true  "Novo horário para o registro"
// @Success      200      {object}  map[string]string
// @Failure      400      {string}  string  "Invalid ID format or request body"
// @Failure      404      {string}  string  "Ponto not found or permission denied"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /pontos/{id} [put]
func AtualizarPonto(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user ID from context")
		return
	}

	idParam := chi.URLParam(r, "id")
	pontoID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var payload PontoUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	res, err := database.DB.Exec(
		"UPDATE pontos SET horario = $1 WHERE id = $2 AND user_id = $3",
		payload.Horario, pontoID, userID,
	)
	if err != nil {
		log.Printf("Error updating 'ponto': %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update 'ponto'")
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to confirm update")
		return
	}

	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "'Ponto' not found or you don't have permission to update it")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Ponto updated successfully"})
}

// DeletarPonto godoc
// @Summary      Deleta um registro de ponto
// @Description  Deleta um registro de ponto existente.
// @Tags         Pontos
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "ID do Ponto"
// @Success      204  {string}  string "No Content"
// @Failure      400  {string}  string  "Invalid ID format"
// @Failure      404  {string}  string  "Ponto not found or permission denied"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /pontos/{id} [delete]
func DeletarPonto(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user ID from context")
		return
	}

	idParam := chi.URLParam(r, "id")
	pontoID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	res, err := database.DB.Exec("DELETE FROM pontos WHERE id = $1 AND user_id = $2", pontoID, userID)
	if err != nil {
		log.Printf("Error deleting 'ponto': %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete 'ponto'")
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to confirm deletion")
		return
	}

	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "'Ponto' not found or you don't have permission to delete it")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
