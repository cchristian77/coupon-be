package post

import (
	"base_project/domain/enums"
	sharedErrs "base_project/shared/errors"
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
)

func (b *base) UpdateStatus(ctx context.Context, id uint64, status enums.PostStatus) error {
	logger.Info(ctx, "Update Post Status: Post ID %d to status %s", id, status.String())

	authUser := constant.AuthUserFromCtx(ctx)
	if authUser == nil {
		return sharedErrs.UnauthorizedErr
	}

	post, err := b.repository.FindPostByIDAndUserID(ctx, id, authUser.ID)
	if err != nil {
		return err
	}

	if post.Status == status {
		return sharedErrs.NewBusinessValidationErr("Update Failed. Post's status is already %s.", status.String())
	}

	if _, err = b.repository.UpdatePostStatus(ctx, id, status); err != nil {
		return err
	}

	return nil
}
