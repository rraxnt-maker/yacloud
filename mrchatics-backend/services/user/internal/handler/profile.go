package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type ProfileHandler struct {
    repo *repository.ProfileRepository
}

func NewProfileHandler(repo *repository.ProfileRepository) *ProfileHandler {
    return &ProfileHandler{repo: repo}
}

// GetProfile godoc
// @Summary Get user profile
// @Tags profile
// @Param user_id path string true "User ID"
// @Success 200 {object} models.UserProfile
// @Router /api/v1/profile/{user_id} [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := uuid.Parse(vars["user_id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    profile, err := h.repo.GetProfile(userID)
    if err != nil {
        http.Error(w, "Profile not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags profile
// @Param user_id path string true "User ID"
// @Param request body models.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} map[string]string
// @Router /api/v1/profile/{user_id} [put]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := uuid.Parse(vars["user_id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    var req models.UpdateProfileRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    if err := h.repo.UpdateProfile(userID, &req); err != nil {
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
