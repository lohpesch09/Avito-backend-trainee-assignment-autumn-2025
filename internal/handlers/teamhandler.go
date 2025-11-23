package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/errors"
	t "github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/team"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/repositories"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/services"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/store"
)

type TeamHandler struct {
	store *store.Store
}

func NewTeamHandler(s *store.Store) *TeamHandler {
	return &TeamHandler{
		store: s,
	}
}

func (h *TeamHandler) TeamCreateHandler(w http.ResponseWriter, r *http.Request) {
	teamRepo := repositories.NewTeamRepository(h.store)
	userRepo := repositories.NewUserRepository(h.store)
	teamService := services.NewTeamService(teamRepo, userRepo)
	team := &t.Team{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(team); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	if err := team.Validation(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	teamResponse, err := teamService.TeamCreate(team)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]t.Team{
		"team": *teamResponse,
	})
}

func (h *TeamHandler) TeamGetHandler(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	teamRepo := repositories.NewTeamRepository(h.store)
	userRepo := repositories.NewUserRepository(h.store)
	teamService := services.NewTeamService(teamRepo, userRepo)
	team, err := teamService.TeamGet(teamName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	} else if team == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: errors.NotFound})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}
