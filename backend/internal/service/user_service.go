package service

import (
	"fmt"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) interfaces.UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetProfile(userID string) (*entity.User, *entity.UserStats, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %w", err)
	}

	stats, err := s.userRepo.GetStats(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return user, stats, nil
}

func (s *UserService) GetContacts(userID string) ([]*entity.User, error) {
	// В реальности нужно получать контакты из Max API
	// Пока возвращаем пустой список
	return []*entity.User{}, nil
}
