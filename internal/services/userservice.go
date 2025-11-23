package services

import (
	"database/sql"

	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/errors"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/pullrequest"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/user"
	"github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) UserSetIsActive(userId string, isActive bool) (*user.User, error) {
	_, err := s.userRepo.FindUserById(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound
		}
		return nil, err
	}
	if err := s.userRepo.SetIsActive(userId, isActive); err != nil {
		return nil, err
	}
	u, err := s.userRepo.FindUserById(userId)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) GetUserReviews(userId string) ([]*pullrequest.PullRequestShort, error) {
	return s.userRepo.FindReviewsByUserId(userId)
}
