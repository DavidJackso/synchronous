package service

import (
	"fmt"
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
	// В реальности нужно валидировать maxToken через Max API
	// Для примера просто создаем пользователя

	// Получаем информацию о боте (для проверки токена)
	_, err := s.maxAPISvc.GetBotInfo()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to validate max token: %w", err)
	}

	// Здесь должна быть логика получения информации о пользователе из Max API
	// Пока используем заглушку
	maxUserID := int64(123456789) // В реальности получаем из Max API

	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByMaxUserID(maxUserID)
	if err != nil {
		// Создаем нового пользователя
		user = &entity.User{
			ID:        uuid.New().String(),
			Name:      "User", // В реальности получаем из Max API
			MaxUserID: maxUserID,
			CreatedAt: time.Now(),
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
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

	// Получаем TTL из конфигурации (нужно передать через конструктор)
	accessTTL := 3600 * time.Second // По умолчанию 1 час

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

	// Получаем TTL из конфигурации
	accessTTL := 3600 * time.Second // По умолчанию 1 час

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
