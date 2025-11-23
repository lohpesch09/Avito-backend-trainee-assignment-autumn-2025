package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/pullrequest"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/store"
)

type PullRequestRepository struct {
	Store *store.Store
}

func NewPullRequestRepository(s *store.Store) *PullRequestRepository {
	return &PullRequestRepository{
		Store: s,
	}
}

func (r *PullRequestRepository) CreatePullRequestWithTx(pr *pullrequest.PullRequest, tx *sql.Tx) error {
	if _, err := tx.Exec("INSERT INTO pull_requests VALUES ($1, $2, $3, $4, $5, $6);",
		pr.PullRequestShort.PullRequestId,
		pr.PullRequestShort.PullRequestName,
		pr.PullRequestShort.AuthorId,
		pr.PullRequestShort.Status,
		pr.CreatedAt,
		pr.MergedAt,
	); err != nil {
		return err
	}
	return nil
}

func (r *PullRequestRepository) CreatePullRequestReviewerWithTx(pullRequestId string, reviewersId []string, tx *sql.Tx) error {
	for _, reviewerId := range reviewersId {
		if _, err := tx.Exec("INSERT INTO pr_reviewers VALUES ($1, $2);",
			pullRequestId,
			reviewerId,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *PullRequestRepository) FindPullRequestById(pullRequestId string, tx *sql.Tx) (*pullrequest.PullRequest, error) {
	pullRequestShort := &pullrequest.PullRequestShort{}
	pullRequest := pullrequest.NewPullRequest(pullRequestShort)
	pullRequestQuery := &sql.Row{}
	if tx != nil {
		pullRequestQuery = tx.QueryRow("SELECT pr_id, pr_name, author_id, status, merged_at "+
			"FROM pull_requests WHERE pr_id = $1;",
			pullRequestId,
		)
	} else {
		pullRequestQuery = r.Store.DB.QueryRow("SELECT pr_id, pr_name, author_id, status, merged_at "+
			"FROM pull_requests WHERE pr_id = $1;",
			pullRequestId,
		)
	}
	err := pullRequestQuery.Scan(
		&pullRequestShort.PullRequestId,
		&pullRequestShort.PullRequestName,
		&pullRequestShort.AuthorId,
		&pullRequestShort.Status,
		&pullRequest.MergedAt,
	)
	if err != nil {
		return nil, err
	}
	var reviewId string
	reviewersQuery := &sql.Rows{}
	if tx != nil {
		reviewersQuery, err = tx.Query("SELECT user_id FROM pr_reviewers WHERE pr_id = $1;",
			pullRequestId,
		)
		if err != nil {
			return nil, err
		}
	} else {
		reviewersQuery, err = r.Store.DB.Query("SELECT user_id FROM pr_reviewers WHERE pr_id = $1;",
			pullRequestId,
		)
		if err != nil {
			return nil, err
		}
	}
	for reviewersQuery.Next() {
		if err := reviewersQuery.Scan(&reviewId); err != nil {
			if err == sql.ErrNoRows {
				return pullRequest, nil
			}
			return nil, err
		}
		pullRequest.AssignedReviwers = append(pullRequest.AssignedReviwers, reviewId)
	}
	return pullRequest, nil
}

func (r *PullRequestRepository) MergePullRequest(pullRequestId string, tx *sql.Tx) error {
	if _, err := tx.Exec("UPDATE pull_requests SET status = $1, merged_at = NOW() WHERE pr_id = $2;",
		pullrequest.MERGED,
		pullRequestId,
	); err != nil {
		return err
	}
	return nil
}

func (r *PullRequestRepository) CheckIfPullRequestMerged(pullRequestId string, tx *sql.Tx) (bool, error) {
	var isMerged bool
	isMergedRow := tx.QueryRow("SELECT status = $1 FROM pull_requests WHERE pr_id = $2;",
		pullrequest.MERGED,
		pullRequestId,
	)
	if err := isMergedRow.Scan(&isMerged); err != nil {
		return isMerged, err
	}
	return isMerged, nil
}

func (r *PullRequestRepository) CheckIfUserIsReviewer(pullRequestId string, oldUserId string, tx *sql.Tx) (bool, error) {
	var isReviewer bool
	isReviewerRow := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2);",
		pullRequestId,
		oldUserId,
	)
	if err := isReviewerRow.Scan(&isReviewer); err != nil {
		return isReviewer, err
	}
	return isReviewer, nil
}

func (r *PullRequestRepository) ReassignReviewer(pullRequestId string, oldUserId string, tx *sql.Tx) (string, error) {
	var oldUserTeamName string
	var newReviewerId string
	var authorId string
	oldUserTeamNameRow := tx.QueryRow("SELECT team_name FROM users WHERE user_id = $1;", oldUserId)
	if err := oldUserTeamNameRow.Scan(&oldUserTeamName); err != nil {
		return "", err
	}
	authorIdRow := tx.QueryRow("SELECT author_id FROM pull_requests WHERE pr_id = $1;", pullRequestId)
	if err := authorIdRow.Scan(&authorId); err != nil {
		return "", err
	}
	fmt.Println(oldUserId, oldUserTeamName, authorId)
	newReviewerIdRow := tx.QueryRow("SELECT user_id FROM users "+
		"WHERE team_name = $1 AND is_active = true AND user_id NOT IN "+
		"(SELECT user_id FROM pr_reviewers WHERE pr_id = $2) AND user_id != $3;",
		oldUserTeamName,
		pullRequestId,
		authorId,
	)
	if err := newReviewerIdRow.Scan(&newReviewerId); err != nil {
		return "", err
	}
	fmt.Println(newReviewerId)
	if _, err := tx.Exec("DELETE FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2;", pullRequestId, oldUserId); err != nil {
		return "", err
	}
	if _, err := tx.Exec("INSERT INTO pr_reviewers VALUES ($1, $2);", pullRequestId, newReviewerId); err != nil {
		return "", err
	}
	return newReviewerId, nil
}
