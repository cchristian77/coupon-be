package auth

import (
	sharedErrs "base_project/shared/errors"
	"base_project/util/constant"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (b *base) Logout(ctx context.Context) error {
	sessionID, ok := ctx.Value(constant.XSessionIDKey).(string)
	if !ok {
		return sharedErrs.UnauthorizedErr
	}

	session, err := b.repository.FindSessionBySessionID(ctx, sessionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if session == nil {
		return sharedErrs.InvalidTokenErr
	}

	if err = b.repository.DeleteSessionByID(ctx, session.ID); err != nil {
		return err
	}

	return nil
}
