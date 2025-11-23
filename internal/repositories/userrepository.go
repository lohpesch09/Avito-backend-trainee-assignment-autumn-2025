package repositories

import (
	"database/sql"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/pullrequest"
	t "github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/team"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/user"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/store"
)

type UserRepository struct {
	Store *store.Store
}

func NewUserRepository(s *store.Store) *UserRepository {
	return &UserRepository{
		Store: s,
	}
}

func (r *UserRepository) FindUserById(userId string) (*user.User, error) {
	member := &t.TeamMember{}
	u := &user.User{}
	u.TeamMember = member
	userIdQuery := r.Store.DB.QueryRow("SELECT user_id, user_name, team_name, is_active "+
		"FROM users WHERE user_id = $1;",
		userId,
	)
	if err := userIdQuery.Scan(
		&u.TeamMember.UserId,
		&u.TeamMember.UserName,
		&u.TeamName,
		&u.TeamMember.IsActive,
	); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) CreateUsersWithTx(members []t.TeamMember, teamName string, tx *sql.Tx) error {
	for _, member := range members {
		_, err := tx.Exec("INSERT INTO users VALUES ($1, $2, $3, $4);",
			member.UserId,
			member.UserName,
			teamName,
			member.IsActive,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) SetIsActive(userId string, isActive bool) error {
	if _, err := r.Store.DB.Exec("UPDATE users SET is_active = $1 WHERE user_id = $2;",
		isActive,
		userId,
	); err != nil {
		return err
	}
	if !isActive {
		if _, err := r.Store.DB.Exec("DELETE FROM pr_reviewers WHERE user_id = $1;", userId); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) FindReviewers(userId string) ([]string, error) {
	var teamAuthorName string
	reviewers := make([]string, 0, 2)
	var reviewerId string
	teamNameRow := r.Store.DB.QueryRow("SELECT team_name FROM users WHERE user_id = $1;", userId)
	if err := teamNameRow.Scan(&teamAuthorName); err != nil {
		return nil, err
	}
	reviewersRows, err := r.Store.DB.Query("SELECT user_id FROM users WHERE team_name = $1 AND user_id != $2 "+
		"AND is_active = true LIMIT 2;",
		teamAuthorName,
		userId,
	)
	if err != nil {
		return nil, err
	}
	for reviewersRows.Next() {
		if err := reviewersRows.Scan(&reviewerId); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewerId)
	}
	return reviewers, nil
}

func (r *UserRepository) FindReviewsByUserId(userId string) ([]*pullrequest.PullRequestShort, error) {
	pullRequestShort := pullrequest.PullRequestShort{}
	pullRequestShortSlice := make([]*pullrequest.PullRequestShort, 0)
	reviewsRows, err := r.Store.DB.Query("SELECT pr.pr_id, pr.pr_name, pr.author_id, pr.status FROM pull_requests pr "+
		"INNER JOIN pr_reviewers prr ON pr.pr_id = prr.pr_id AND prr.user_id = $1;",
		userId,
	)
	if err != nil {
		return nil, err
	}
	for reviewsRows.Next() {
		if err := reviewsRows.Scan(
			&pullRequestShort.PullRequestId,
			&pullRequestShort.PullRequestName,
			&pullRequestShort.AuthorId,
			&pullRequestShort.Status,
		); err != nil {
			if err == sql.ErrNoRows {
				return pullRequestShortSlice, nil
			}
			return nil, err
		}
		pullRequestShortSlice = append(pullRequestShortSlice, &pullRequestShort)
	}
	return pullRequestShortSlice, nil
}
