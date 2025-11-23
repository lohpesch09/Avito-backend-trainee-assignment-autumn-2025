package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/errors"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/repositories"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/services"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/store"
)

type UserHandler struct {
	Store *store.Store
}

type UserActiveUpdate struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func NewUserHandler(s *store.Store) *UserHandler {
	return &UserHandler{
		Store: s,
	}
}

func (h *UserHandler) UserSetIsActiveHandler(w http.ResponseWriter, r *http.Request) {
	userRepo := *repositories.NewUserRepository(h.Store)
	userService := *services.NewUserService(&userRepo)
	userActiveUpdate := &UserActiveUpdate{}
	if err := json.NewDecoder(r.Body).Decode(userActiveUpdate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	u, err := userService.UserSetIsActive(userActiveUpdate.UserId, userActiveUpdate.IsActive)
	if err != nil {
		if err == errors.NotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":   u.TeamMember.UserId,
		"username":  u.TeamMember.UserName,
		"team_name": u.TeamName,
		"is_active": u.TeamMember.IsActive,
	})
}

func (h *UserHandler) UserGetReviewHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("user_id")
	userRepo := *repositories.NewUserRepository(h.Store)
	userService := *services.NewUserService(&userRepo)
	reviews, err := userService.GetUserReviews(userName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       userName,
		"pull_requests": reviews,
	})
}
