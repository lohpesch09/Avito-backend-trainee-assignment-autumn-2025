package main

import (
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/server"
)

func main() {
	c := config.NewConfig()
	s := server.NewServer(c)
	err := s.Start()
	if err != nil {
		s.Logger.Error(err.Error())
	}
}
