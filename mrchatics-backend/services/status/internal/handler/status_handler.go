package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"status-service/internal/model"
	"status-service/internal/service"
)

type StatusHandler struct {
	service *service.StatusService
}

func NewStatusHandler(service *service.StatusService) *StatusHandler {
	return &StatusHandler{service: service}
}

func extractUserID(path string) (int, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return 0, http.ErrNoCookie
	}
	userIDStr := parts[3]
	return strconv.Atoi(userIDStr)
}

func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	userID, err := extractUserID(r.URL.Path)
	if err != nil {
		h.writeError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	status, err := h.service.GetUserStatus(userID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	
	h.writeJSON(w, status, http.StatusOK)
}

func (h *StatusHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req model.UpdateStatusRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.service.UpdateStatus(&req); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	h.writeJSON(w, map[string]string{"message": "Status updated successfully"}, http.StatusOK)
}

func (h *StatusHandler) GetBatchStatuses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req model.BatchStatusRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	statuses, err := h.service.GetBatchStatuses(req.UserIDs)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := model.BatchStatusResponse{Statuses: statuses}
	h.writeJSON(w, response, http.StatusOK)
}

func (h *StatusHandler) GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	users, err := h.service.GetOnlineUsers()
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	h.writeJSON(w, map[string][]int{"online_users": users}, http.StatusOK)
}

func (h *StatusHandler) UpdateLastSeen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		UserID int `json:"user_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.service.UpdateLastSeen(req.UserID); err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	h.writeJSON(w, map[string]string{"message": "Last seen updated"}, http.StatusOK)
}

func (h *StatusHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	h.writeJSON(w, map[string]string{"status": "healthy", "service": "status-service"}, http.StatusOK)
}

func (h *StatusHandler) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *StatusHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.ErrorResponse{Error: message})
}
