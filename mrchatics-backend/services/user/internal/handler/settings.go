package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type SettingsHandler struct {
    repo *repository.ProfileRepository
}

func NewSettingsHandler(repo *repository.ProfileRepository) *SettingsHandler {
    return &SettingsHandler{repo: repo}
}

// GetSettings godoc
// @Summary Get user settings
// @Tags settings
// @Param user_id path string true "User ID"
// @Success 200 {object} models.UserSettings
// @Router /api/v1/settings/{user_id} [get]
func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := uuid.Parse(vars["user_id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    settings, err := h.repo.GetSettings(userID)
    if err != nil {
        http.Error(w, "Settings not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(settings)
}

// UpdateSettings godoc
// @Summary Update user settings
// @Tags settings
// @Param user_id path string true "User ID"
// @Param request body models.UpdateSettingsRequest true "Update settings request"
// @Success 200 {object} map[string]string
// @Router /api/v1/settings/{user_id} [put]
func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := uuid.Parse(vars["user_id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    var req models.UpdateSettingsRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    if err := h.repo.UpdateSettings(userID, &req); err != nil {
        http.Error(w, "Failed to update settings", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Settings updated successfully"})
}

// BatchUpdate godoc
// @Summary Batch update profile and settings
// @Tags batch
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/batch/{user_id} [put]
func (h *SettingsHandler) BatchUpdate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := uuid.Parse(vars["user_id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    var updates struct {
        Profile  *models.UpdateProfileRequest  `json:"profile"`
        Settings *models.UpdateSettingsRequest `json:"settings"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Update both in transaction
    if updates.Profile != nil {
        if err := h.repo.UpdateProfile(userID, updates.Profile); err != nil {
            http.Error(w, "Failed to update profile", http.StatusInternalServerError)
            return
        }
    }
    
    if updates.Settings != nil {
        if err := h.repo.UpdateSettings(userID, updates.Settings); err != nil {
            http.Error(w, "Failed to update settings", http.StatusInternalServerError)
            return
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Batch update completed successfully"})
}
