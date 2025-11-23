package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/errors"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/pullrequest"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/repositories"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/services"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/store"
)

type PullRequestHandler struct {
	Store *store.Store
}

type PullRequestReassign struct {
	PullRequestId string `json:"pull_request_id"`
	OldUserId     string `json:"old_user_id"`
}

func NewPullRequestHandler(s *store.Store) *PullRequestHandler {
	return &PullRequestHandler{
		Store: s,
	}
}

func (h *PullRequestHandler) PullRequestCreateHandler(w http.ResponseWriter, r *http.Request) {
	userRepo := repositories.NewUserRepository(h.Store)
	pullRequestRepo := repositories.NewPullRequestRepository(h.Store)
	pullRequestService := services.NewPullRequestService(userRepo, pullRequestRepo)
	pullRequestShort := &pullrequest.PullRequestShort{}
	if err := json.NewDecoder(r.Body).Decode(pullRequestShort); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	if err := pullRequestShort.Validation(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	pullRequestShort.Status = pullrequest.OPEN
	pullRequest := &pullrequest.PullRequest{}
	pullRequest.PullRequestShort = pullRequestShort
	pullRequest.CreatedAt = time.Now()
	pullRequestResponse, err := pullRequestService.PullRequestCreate(pullRequest)
	if err != nil {
		if err == errors.NotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
			return
		}
		if err == errors.PrExists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pull_request_id":    pullRequestResponse.PullRequestShort.PullRequestId,
		"pull_request_name":  pullRequestResponse.PullRequestShort.PullRequestName,
		"author_id":          pullRequestResponse.PullRequestShort.AuthorId,
		"status":             pullRequestResponse.PullRequestShort.Status,
		"assigned_reviewers": pullRequestResponse.AssignedReviwers,
	})
}

func (h *PullRequestHandler) PullRequestMergeHandler(w http.ResponseWriter, r *http.Request) {
	userRepo := repositories.NewUserRepository(h.Store)
	pullRequestRepo := repositories.NewPullRequestRepository(h.Store)
	pullRequestService := services.NewPullRequestService(userRepo, pullRequestRepo)
	pullRequestShort := &pullrequest.PullRequestShort{}
	if err := json.NewDecoder(r.Body).Decode(pullRequestShort); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	pullRequestResponse, err := pullRequestService.PullRequestMerge(pullRequestShort.PullRequestId)
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
		"pull_request_id":    pullRequestResponse.PullRequestShort.PullRequestId,
		"pull_request_name":  pullRequestResponse.PullRequestShort.PullRequestName,
		"author_id":          pullRequestResponse.PullRequestShort.AuthorId,
		"status":             pullRequestResponse.PullRequestShort.Status,
		"assigned_reviewers": pullRequestResponse.AssignedReviwers,
		"mergedAt":           pullRequestResponse.MergedAt,
	})
}

func (h *PullRequestHandler) PullRequestReassignHandler(w http.ResponseWriter, r *http.Request) {
	userRepo := repositories.NewUserRepository(h.Store)
	pullRequestRepo := repositories.NewPullRequestRepository(h.Store)
	pullRequestService := services.NewPullRequestService(userRepo, pullRequestRepo)
	pullRequestReassign := &PullRequestReassign{}
	if err := json.NewDecoder(r.Body).Decode(pullRequestReassign); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	pullRequestResponse, replacedByUserId, err := pullRequestService.PullRequestReassign(pullRequestReassign.PullRequestId,
		pullRequestReassign.OldUserId,
	)
	if err != nil {
		if err == errors.NotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
			return
		}
		if err == errors.PrMerged || err == errors.NotAssigned || err == errors.NoCandidate {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": map[string]interface{}{
			"pull_request_id":    pullRequestResponse.PullRequestShort.PullRequestId,
			"pull_request_name":  pullRequestResponse.PullRequestShort.PullRequestName,
			"author_id":          pullRequestResponse.PullRequestShort.AuthorId,
			"status":             pullRequestResponse.PullRequestShort.Status,
			"assigned_reviewers": pullRequestResponse.AssignedReviwers,
		},
		"replaced_by": replacedByUserId,
	})
}
