package repositories

import (
	"database/sql"
	"errors"

	t "github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/team"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/store"
)

type TeamRepository struct {
	Store *store.Store
}

func NewTeamRepository(s *store.Store) *TeamRepository {
	return &TeamRepository{
		Store: s,
	}
}

func (r *TeamRepository) FindTeamByName(teamName string) error {
	var team string
	teamQuery := r.Store.DB.QueryRow("SELECT * FROM teams WHERE team_name = $1;", teamName)
	if err := teamQuery.Scan(&team); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return errors.New("team exists")
}

func (r *TeamRepository) CreateWithTx(teamName string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO teams values ($1)", teamName)
	return err
}

func (r *TeamRepository) FindTeamWithMembersByName(teamName string) (*t.Team, error) {
	team := &t.Team{}
	members, err := r.Store.DB.Query("SELECT user_id, user_name, is_active "+
		"FROM users WHERE team_name = $1;",
		teamName,
	)
	if err != nil {
		return team, err
	}
	for members.Next() {
		teamMember := &t.TeamMember{}
		if err := members.Scan(&teamMember.UserId,
			&teamMember.UserName,
			&teamMember.IsActive,
		); err != nil {
			return &t.Team{}, err
		}
		team.Members = append(team.Members, *teamMember)
	}
	if len(team.Members) != 0 {
		team.TeamName = teamName
	} else {
		return nil, nil
	}
	return team, nil
}
