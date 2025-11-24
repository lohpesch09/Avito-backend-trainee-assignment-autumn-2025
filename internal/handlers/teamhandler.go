package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/errors"
	t "github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/team"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/services"
)

type TeamHandler struct {
	TeamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{
		TeamService: teamService,
	}
}

func (h *TeamHandler) TeamCreateHandler(w http.ResponseWriter, r *http.Request) {
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
	teamResponse, err := h.TeamService.TeamCreate(team)
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
	team, err := h.TeamService.TeamGet(teamName)
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
