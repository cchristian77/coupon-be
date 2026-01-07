package auth

import (
	"base_project/domain"
	sharedErrs "base_project/shared/errors"
	"base_project/util/logger"
	tokenMaker "base_project/util/token"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Authenticate functions to decrypt and verify the provided access token.
func (b *base) Authenticate(ctx context.Context, accessToken string) (*domain.User, *tokenMaker.Payload, error) {
	logger.Info(ctx, "Authentication with access token: %v", accessToken)

	payload, err := tokenMaker.Get().Verify(accessToken)
	if err != nil {
		return nil, nil, err
	}

	// find session by verified payload ID
	session, err := b.repository.FindSessionBySessionID(ctx, payload.ID.String())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	if session == nil {
		return nil, nil, sharedErrs.InvalidTokenErr
	}

	if session.AccessToken != accessToken {
		return nil, nil, sharedErrs.InvalidTokenErr
	}

	// check whether session is expired
	if time.Now().After(session.AccessTokenExpiresAt) {
		return nil, nil, sharedErrs.ExpiredTokenErr
	}

	// find the user data from payload
	authUser, err := b.repository.FindUserByID(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, sharedErrs.IncorrectCredentialErr
		}

		return nil, nil, err
	}

	return authUser, payload, nil
}
