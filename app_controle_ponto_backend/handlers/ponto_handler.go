package handlers

import (
	"bytes"
	"controle-ponto-api/database"
	"controle-ponto-api/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// --- Funções Auxiliares ---

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"erro": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func parseDateParam(r *http.Request) (string, error) {
	dataParam := chi.URLParam(r, "data")
	if _, err := time.Parse("2006-01-02", dataParam); err != nil {
		return "", fmt.Errorf("formato de data inválido. Use AAAA-MM-DD")
	}
	return dataParam, nil
}

// --- Handlers ---

// PontoRequest representa o corpo da requisição para registrar um ponto.
type PontoRequest struct {
	Data   string `json:"data"` // Formato: AAAA-MM-DD
	Hora   int    `json:"hora"`
	Minuto int    `json:"minuto"`
}

// RegistrarPonto é o handler para a rota que registra um novo ponto.
func RegistrarPonto(w http.ResponseWriter, r *http.Request) {
	// Read the request body into a byte slice for logging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Erro ao ler corpo da requisição: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}
	// Restore the body for subsequent reads
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	log.Printf("Requisição recebida para /registrar-ponto. Corpo: %s", string(bodyBytes))

	var req PontoRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Erro ao decodificar corpo da requisição: %v", err)
		respondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	log.Printf("Dados decodificados: Data=%s, Hora=%d, Minuto=%d", req.Data, req.Hora, req.Minuto)

	var validationErrors []string
	var data time.Time

	// Valida a data
	data, err = time.Parse("2006-01-02", req.Data)
	if err != nil {
		validationErrors = append(validationErrors, "Formato de data inválido. Use AAAA-MM-DD.")
	}

	// Valida a hora
	if req.Hora < 0 || req.Hora > 23 {
		validationErrors = append(validationErrors, "Valor de hora inválido. Use um valor entre 0 e 23.")
	}

	// Valida o minuto
	if req.Minuto < 0 || req.Minuto > 59 {
		validationErrors = append(validationErrors, "Valor de minuto inválido. Use um valor entre 0 e 59.")
	}

	// Se houver erros de validação, retorna todos eles
	if len(validationErrors) > 0 {
		respondWithJSON(w, http.StatusBadRequest, map[string][]string{"erros": validationErrors})
		return
	}

	// Se a validação passou, prossegue com a lógica de negócio
	horarioDoPonto := time.Date(data.Year(), data.Month(), data.Day(), req.Hora, req.Minuto, 0, 0, time.Local)

	// Verifica se já existe um ponto com o mesmo horário
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM pontos WHERE horario = ?", horarioDoPonto).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar duplicidade de ponto: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor ao verificar duplicidade")
		return
	}

	if count > 0 {
		respondWithError(w, http.StatusConflict, "Ponto com este horário já registrado.")
		return
	}

	stmt, err := database.DB.Prepare("INSERT INTO pontos(horario) VALUES(?)")
	if err != nil {
		log.Printf("Erro ao preparar statement de inserção: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(horarioDoPonto)
	if err != nil {
		log.Printf("Erro ao executar inserção: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Erro ao obter o último ID inserido: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}

	novoPonto := models.Ponto{
		ID:      strconv.FormatInt(id, 10),
		Horario: horarioDoPonto,
	}

	respondWithJSON(w, http.StatusCreated, novoPonto)
}

// ListarPontosPorData é o handler para a rota que lista os pontos de uma data.
func ListarPontosPorData(w http.ResponseWriter, r *http.Request) {
	dataParam, err := parseDateParam(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	query := `SELECT id, horario FROM pontos WHERE date(horario, 'localtime') = ? ORDER BY horario ASC`
	rows, err := database.DB.Query(query, dataParam)
	if err != nil {
		log.Printf("Erro ao consultar pontos por data: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}
	defer rows.Close()

	pontos := []models.Ponto{}
	for rows.Next() {
		var p models.Ponto
		var id int64
		if err := rows.Scan(&id, &p.Horario); err != nil {
			log.Printf("Erro ao escanear linha do banco: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
			return
		}
		p.ID = strconv.FormatInt(id, 10)
		pontos = append(pontos, p)
	}

	if len(pontos) == 0 {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	respondWithJSON(w, http.StatusOK, pontos)
}

// CalcularHorasTrabalhadas é o handler para a rota que calcula o total de horas de um dia.
func CalcularHorasTrabalhadas(w http.ResponseWriter, r *http.Request) {
	dataParam, err := parseDateParam(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	query := `SELECT horario FROM pontos WHERE date(horario, 'localtime') = ? ORDER BY horario ASC`
	rows, err := database.DB.Query(query, dataParam)
	if err != nil {
		log.Printf("Erro ao consultar pontos para cálculo: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}
	defer rows.Close()

	var horarios []time.Time
	for rows.Next() {
		var horario time.Time
		if err := rows.Scan(&horario); err != nil {
			log.Printf("Erro ao escanear linha para cálculo: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
			return
		}
		horarios = append(horarios, horario)
	}

	var totalDuracao time.Duration
	// Se o número de registros for ímpar, o último é ignorado, o que está correto.
	for i := 0; i < len(horarios)-1; i += 2 {
		entrada := horarios[i]
		saida := horarios[i+1]
		totalDuracao += saida.Sub(entrada)
	}

	totalHoras := int(totalDuracao.Hours())
	totalMinutos := int(totalDuracao.Minutes()) % 60
	totalSegundos := int(totalDuracao.Seconds()) % 60

	resposta := map[string]string{
		"total_trabalhado": fmt.Sprintf("%dh %dm %ds", totalHoras, totalMinutos, totalSegundos),
		"total_segundos":   strconv.FormatFloat(totalDuracao.Seconds(), 'f', 0, 64),
	}

	respondWithJSON(w, http.StatusOK, resposta)
}
