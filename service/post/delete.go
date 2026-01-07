package post

import (
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
)

func (b *base) Delete(ctx context.Context, id uint64) error {
	logger.Info(ctx, "Delete Post with Post ID: %d", id)

	authUser := constant.AuthUserFromCtx(ctx)

	_, err := b.repository.FindPostByIDAndUserID(ctx, id, authUser.ID)
	if err != nil {
		return err
	}

	if err = b.repository.DeletePostByID(ctx, id); err != nil {
		return err
	}

	return nil
}
