package interfaces

import (
	"github.com/rnegic/synchronous/internal/entity"
)

type UserService interface {
	GetProfile(userID string) (*entity.User, *entity.UserStats, error)
	GetContacts(userID string) ([]*entity.User, error)
}
