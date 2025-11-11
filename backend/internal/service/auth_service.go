package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
	"github.com/rnegic/synchronous/pkg/jwt"
)

type AuthService struct {
	userRepo     interfaces.UserRepository
	maxAPISvc    interfaces.MaxAPIService
	tokenManager *jwt.TokenManager
}

func NewAuthService(
	userRepo interfaces.UserRepository,
	maxAPISvc interfaces.MaxAPIService,
	tokenManager *jwt.TokenManager,
) interfaces.AuthService {
	return &AuthService{
		userRepo:     userRepo,
		maxAPISvc:    maxAPISvc,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Login(maxToken, deviceID string) (*entity.AuthTokens, *entity.User, error) {
	profile, err := s.maxAPISvc.GetProfileByToken(maxToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to validate max token: %w", err)
	}

	if profile == nil || profile.UserID == 0 {
		return nil, nil, fmt.Errorf("invalid max token: missing profile information")
	}

	displayName := strings.TrimSpace(fmt.Sprintf("%s %s", profile.FirstName, profile.LastName))
	if displayName == "" {
		displayName = strings.TrimSpace(profile.Name)
	}
	if displayName == "" {
		displayName = fmt.Sprintf("user-%d", profile.UserID)
	}

	var avatarURL *string
	if strings.TrimSpace(profile.AvatarURL) != "" {
		avatar := strings.TrimSpace(profile.AvatarURL)
		avatarURL = &avatar
	}

	user, err := s.userRepo.GetByMaxUserID(profile.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	now := time.Now()

	if user == nil {
		user = &entity.User{
			ID:        uuid.New().String(),
			Name:      displayName,
			AvatarURL: avatarURL,
			MaxUserID: profile.UserID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		needsUpdate := false

		if user.Name != displayName {
			user.Name = displayName
			needsUpdate = true
		}

		if !equalPointers(user.AvatarURL, avatarURL) {
			user.AvatarURL = avatarURL
			needsUpdate = true
		}

		if needsUpdate {
			user.UpdatedAt = now
			if err := s.userRepo.Update(user); err != nil {
				return nil, nil, fmt.Errorf("failed to update user: %w", err)
			}
		}
	}

	// Генерируем токены
	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	accessTTL := s.tokenManager.AccessTTL()

	tokens := &entity.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(accessTTL),
	}

	return tokens, user, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*entity.AuthTokens, error) {
	claims, err := s.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Проверяем, существует ли пользователь
	_, err = s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Генерируем новые токены
	accessToken, err := s.tokenManager.GenerateAccessToken(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	accessTTL := s.tokenManager.AccessTTL()

	tokens := &entity.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(accessTTL),
	}

	return tokens, nil
}

func (s *AuthService) ValidateToken(token string) (string, error) {
	claims, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	return claims.UserID, nil
}

func (s *AuthService) Logout(userID string) error {
	// В реальности можно добавить blacklist токенов
	// Пока просто возвращаем успех
	return nil
}

func equalPointers(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
