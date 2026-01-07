package comment

import (
	sharedErrs "base_project/shared/errors"
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
)

func (b *base) Delete(ctx context.Context, id uint64) error {
	logger.Info(ctx, "Delete Comment with Comment ID: %d", id)

	authUser := constant.AuthUserFromCtx(ctx)
	if authUser == nil {
		return sharedErrs.UnauthorizedErr
	}

	_, err := b.repository.FindCommentByIDAndUserID(ctx, id, authUser.ID)
	if err != nil {
		return err
	}

	if err = b.repository.DeleteCommentByID(ctx, id); err != nil {
		return err
	}

	return nil
}
