package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/api"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/handlers"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/store"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Config *config.Config
	Logger *logrus.Logger
	Router *mux.Router
	Store  *store.Store
	API    *api.API
}

func NewServer(config *config.Config, store *store.Store, API *api.API) *Server {
	return &Server{
		Config: config,
		Logger: logrus.New(),
		Router: mux.NewRouter(),
		Store:  store,
		API:    API,
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()
	time.Sleep(10 * time.Second)
	if err := s.configureStore(); err != nil {
		return err
	}
	s.Logger.Info("server started")
	return http.ListenAndServe(":"+s.Config.Port, s.Router)
}

func (s *Server) configureStore() error {
	db, err := sql.Open("postgres", s.Config.DBURL)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	s.Store.DB = db
	return nil
}

func (s *Server) configureRouter() {
	s.Router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	s.Router.HandleFunc("/team/add", s.API.TeamHandler.TeamCreateHandler).Methods("POST")
	s.Router.HandleFunc("/team/get", s.API.TeamHandler.TeamGetHandler).Methods("GET")
	s.Router.HandleFunc("/users/setIsActive", s.API.UserHandler.UserSetIsActiveHandler).Methods("POST")
	s.Router.HandleFunc("/users/getReview", s.API.UserHandler.UserGetReviewHandler).Methods("GET")
	s.Router.HandleFunc("/pullRequest/create", s.API.PullRequestHandler.PullRequestCreateHandler).Methods("POST")
	s.Router.HandleFunc("/pullRequest/merge", s.API.PullRequestHandler.PullRequestMergeHandler).Methods("POST")
	s.Router.HandleFunc("/pullRequest/reassign", s.API.PullRequestHandler.PullRequestReassignHandler).Methods("POST")
}

func (s *Server) configureLogger() error {
	lvl, err := logrus.ParseLevel(s.Config.LogLevel)
	if err != nil {
		return err
	}
	s.Logger.Level = lvl
	return nil
}
