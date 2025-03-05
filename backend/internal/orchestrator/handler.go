package orchestrator

import (
	"calc-website/internal/models"
	"calc-website/pkg/utils"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type APIHandler struct {
	Service *APIService
}

func NewAPIHandler(service *APIService) *APIHandler {
	return &APIHandler{
		Service: service,
	}
}

func (h *APIHandler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", h.Calculate)
	mux.HandleFunc("/api/v1/expressions", h.GetExpressions)
	mux.HandleFunc("/api/v1/expressions/{id}", h.GetExpressionByID)
	mux.HandleFunc("/internal/task", h.TaskHandler)

	return mux
}

func (h *APIHandler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTask(w, r)
	case http.MethodPost:
		h.PostTask(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *APIHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	task := h.Service.GetTask()
	if task == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(map[string]any{"message": "tasks not found"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h *APIHandler) PostTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer utils.CloseResponseBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var result models.TaskResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	taskID, err := strconv.ParseUint(result.TaskID, 10, 32)
	err = h.Service.ConfirmTask(uint32(taskID), result.Result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *APIHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer utils.CloseResponseBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	var expression models.ExpressionRequest
	err = json.Unmarshal(body, &expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	expressionID, err := h.Service.CreateTasks(expression.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(
		map[string]any{"expression": models.ExpressionResponse{ExpressionID: strconv.Itoa(int(expressionID))}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
}

func (h *APIHandler) GetExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	expressions := h.Service.GetAllExpressions()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]any{"expressions": expressions})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *APIHandler) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	expression := h.Service.GetExpressionByID(uint32(id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]any{"expression": expression})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
