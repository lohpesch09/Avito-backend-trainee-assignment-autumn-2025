package services

import (
	"database/sql"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/errors"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/pullrequest"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/repositories"
)

type PullRequestService struct {
	UserRepo        *repositories.UserRepository
	PullRequestRepo *repositories.PullRequestRepository
}

func NewPullRequestService(userRepo *repositories.UserRepository,
	pullRequestRepo *repositories.PullRequestRepository) *PullRequestService {
	return &PullRequestService{
		UserRepo:        userRepo,
		PullRequestRepo: pullRequestRepo,
	}
}

func (s *PullRequestService) PullRequestCreate(pr *pullrequest.PullRequest) (*pullrequest.PullRequest, error) {
	if _, err := s.UserRepo.FindUserById(pr.PullRequestShort.AuthorId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound
		}
		return nil, err
	}
	reviewersId, err := s.UserRepo.FindReviewers(pr.PullRequestShort.AuthorId)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviwers = reviewersId
	tx, err := s.PullRequestRepo.Store.BeginTx()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	if _, err := s.PullRequestRepo.FindPullRequestById(pr.PullRequestShort.PullRequestId, tx); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == nil {
		return nil, errors.PrExists
	}
	if err := s.PullRequestRepo.CreatePullRequestWithTx(pr, tx); err != nil {
		return nil, err
	}
	if err := s.PullRequestRepo.CreatePullRequestReviewerWithTx(pr.PullRequestShort.PullRequestId, pr.AssignedReviwers, tx); err != nil {
		return nil, err
	}
	tx.Commit()
	return s.PullRequestRepo.FindPullRequestById(pr.PullRequestShort.PullRequestId, nil)
}

func (s *PullRequestService) PullRequestMerge(pullRequestId string) (*pullrequest.PullRequest, error) {
	if _, err := s.PullRequestRepo.FindPullRequestById(pullRequestId, nil); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound
		}
		return nil, err
	}
	tx, err := s.PullRequestRepo.Store.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	isMerged, err := s.PullRequestRepo.CheckIfPullRequestMerged(pullRequestId, tx)
	if err != nil {
		return nil, err
	} else if isMerged {
		return s.PullRequestRepo.FindPullRequestById(pullRequestId, tx)
	}
	s.PullRequestRepo.MergePullRequest(pullRequestId, tx)
	tx.Commit()
	return s.PullRequestRepo.FindPullRequestById(pullRequestId, nil)
}

func (s *PullRequestService) PullRequestReassign(pullRequestId, oldUserId string) (*pullrequest.PullRequest, string, error) {
	if _, err := s.UserRepo.FindUserById(oldUserId); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.NotFound
		}
		return nil, "", err
	}
	if _, err := s.PullRequestRepo.FindPullRequestById(pullRequestId, nil); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.NotFound
		}
		return nil, "", err
	}
	tx, err := s.PullRequestRepo.Store.BeginTx()
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()
	isMerged, err := s.PullRequestRepo.CheckIfPullRequestMerged(pullRequestId, tx)
	if err != nil {
		return nil, "", err
	} else if isMerged {
		return nil, "", errors.PrMerged
	}
	isReviewer, err := s.PullRequestRepo.CheckIfUserIsReviewer(pullRequestId, oldUserId, tx)
	if err != nil {
		return nil, "", err
	} else if !isReviewer {
		return nil, "", errors.NotAssigned
	}
	newReviewerId, err := s.PullRequestRepo.ReassignReviewer(pullRequestId, oldUserId, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.NoCandidate
		}
		return nil, "", err
	}
	pullRequest, err := s.PullRequestRepo.FindPullRequestById(pullRequestId, tx)
	if err != nil {
		return nil, "", err
	}
	tx.Commit()
	return pullRequest, newReviewerId, nil
}
