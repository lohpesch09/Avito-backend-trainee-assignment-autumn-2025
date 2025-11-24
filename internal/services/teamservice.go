package services

import (
	"database/sql"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/errors"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/repositories"

	t "github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/team"
)

type TeamService struct {
	TeamRepo *repositories.TeamRepository
	UserRepo *repositories.UserRepository
}

func NewTeamService(t *repositories.TeamRepository, u *repositories.UserRepository) *TeamService {
	return &TeamService{
		TeamRepo: t,
		UserRepo: u,
	}
}

func (s *TeamService) TeamCreate(team *t.Team) (*t.Team, error) {
	err := s.TeamRepo.FindTeamByName(team.TeamName)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	} else {
		return nil, errors.TeamExists
	}
	for _, member := range team.Members {
		_, err := s.UserRepo.FindUserById(member.UserId)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return nil, err
		} else {
			return nil, errors.UserExistsAnotherTeam
		}
	}
	tx, err := s.TeamRepo.Store.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if err := s.TeamRepo.CreateWithTx(team.TeamName, tx); err != nil {
		return nil, err
	}
	if err := s.UserRepo.CreateUsersWithTx(team.Members, team.TeamName, tx); err != nil {
		return nil, err
	}
	tx.Commit()
	return s.TeamRepo.FindTeamWithMembersByName(team.TeamName)
}

func (s *TeamService) TeamGet(teamName string) (*t.Team, error) {
	return s.TeamRepo.FindTeamWithMembersByName(teamName)
}
