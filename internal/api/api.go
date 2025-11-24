package api

import (
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/handlers"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/repositories"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/services"
)

type API struct {
	UserHandler        *handlers.UserHandler
	TeamHandler        *handlers.TeamHandler
	PullRequestHandler *handlers.PullRequestHandler
}

func NewAPI() *API {
	userRepo := &repositories.UserRepository{}
	teamRepo := &repositories.TeamRepository{}
	pullRequestRepo := &repositories.PullRequestRepository{}

	userService := services.NewUserService(userRepo)
	teamService := services.NewTeamService(teamRepo, userRepo)
	pullRequestService := services.NewPullRequestService(userRepo, pullRequestRepo)

	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	pullRequestHandler := handlers.NewPullRequestHandler(pullRequestService)

	return &API{
		UserHandler:        userHandler,
		TeamHandler:        teamHandler,
		PullRequestHandler: pullRequestHandler,
	}
}
