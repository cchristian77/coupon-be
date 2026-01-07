package auth

import (
	"base_project/domain"
	"base_project/request"
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util/config"
	"base_project/util/logger"
	tokenMaker "base_project/util/token"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (b *base) Login(ctx context.Context, input *request.Login) (*response.Auth, error) {
	logger.Info(ctx, "Login with request: %v", input)

	authUser, err := b.repository.FindUserByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrs.IncorrectCredentialErr
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(authUser.Password), []byte(input.Password))
	if err != nil {
		return nil, sharedErrs.IncorrectCredentialErr
	}

	sessionID := uuid.New()

	accessTokenDuration, _ := time.ParseDuration(config.Env().Auth.AccessTokenExpiration)

	accessToken, payload, err := tokenMaker.Get().Generate(sessionID, authUser.ID, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	_, err = b.repository.CreateSession(ctx, &domain.Session{
		UserID:               authUser.ID,
		SessionID:            payload.ID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: time.Unix(payload.StandardClaims.ExpiresAt, 0),
		AccessTokenCreatedAt: time.Unix(payload.StandardClaims.IssuedAt, 0),
		UserAgent:            input.UserAgent,
		ClientIP:             input.ClientIP,
	})
	if err != nil {
		return nil, err
	}

	return &response.Auth{
		User: response.User{
			ID:       authUser.ID,
			Username: authUser.Username,
			FullName: authUser.FullName,
			Role:     authUser.Role.String(),
		},
		SessionID:            payload.ID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: payload.ExpiresAt,
	}, nil
}
