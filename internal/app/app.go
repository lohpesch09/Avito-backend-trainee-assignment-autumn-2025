package app

import (
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/api"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/server"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/store"
)

func Start() {
	API := api.NewAPI()
	store := store.NewStore()
	API.UserHandler.UserService.UserRepo.Store = store
	API.TeamHandler.TeamService.TeamRepo.Store = store
	API.PullRequestHandler.PullRequestService.PullRequestRepo.Store = store
	c := config.NewConfig()
	s := server.NewServer(c, store, API)
	err := s.Start()
	if err != nil {
		s.Logger.Error(err.Error())
	}
}
