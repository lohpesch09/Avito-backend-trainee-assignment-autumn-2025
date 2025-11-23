package services

import (
	"database/sql"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/errors"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/repositories"

	t "github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/team"
)

type TeamService struct {
	teamRepo *repositories.TeamRepository
	userRepo *repositories.UserRepository
}

func NewTeamService(t *repositories.TeamRepository, u *repositories.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: t,
		userRepo: u,
	}
}

func (s *TeamService) TeamCreate(team *t.Team) (*t.Team, error) {
	err := s.teamRepo.FindTeamByName(team.TeamName)
	if err != nil {
		if err.Error() == "team exists" {
			return nil, errors.TeamExists
		}
		return nil, err
	}
	for _, member := range team.Members {
		_, err := s.userRepo.FindUserById(member.UserId)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return nil, err
		} else {
			return nil, errors.UserExistsAnotherTeam
		}
	}
	tx, err := s.teamRepo.Store.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if err := s.teamRepo.CreateWithTx(team.TeamName, tx); err != nil {
		return nil, err
	}
	if err := s.userRepo.CreateUsersWithTx(team.Members, team.TeamName, tx); err != nil {
		return nil, err
	}
	tx.Commit()
	return s.teamRepo.FindTeamWithMembersByName(team.TeamName)
}

func (s *TeamService) TeamGet(teamName string) (*t.Team, error) {
	return s.teamRepo.FindTeamWithMembersByName(teamName)
}
