package repository

import (
	"context"
	"coupon_be/domain"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/external/database"
	"coupon_be/util"
	"coupon_be/util/logger"

	"gorm.io/gorm/clause"
)

func (r *repo) FindCommentsPaginatedByPostID(ctx context.Context, postID uint64, p *util.Pagination) ([]*domain.Comment, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	var (
		data  []*domain.Comment
		count int64
	)

	query := db.WithContext(ctx).Model(&data).
		Preload("User").
		Where("post_id = ?", postID)

	if err := query.Count(&count).Error; err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find comment paginated count: %v", err)
		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	p.SetTotal(count)

	err := query.Offset(p.Offset()).
		Limit(p.Limit()).
		Find(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find comment paginated: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) FindCommentsPaginatedByUserID(ctx context.Context, userID uint64, p *util.Pagination) ([]*domain.Comment, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	var (
		data  []*domain.Comment
		count int64
	)

	query := db.WithContext(ctx).Model(&data).
		Preload("User").
		Where("user_id = ?", userID)

	if err := query.Count(&count).Error; err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find comment paginated count: %v", err)
		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	p.SetTotal(count)

	err := query.Offset(p.Offset()).
		Limit(p.Limit()).
		Find(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find comment paginated: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) FindCommentByID(ctx context.Context, id uint64) (*domain.Comment, error) {
	var result *domain.Comment

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		First(&result, id).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find comment by id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindCommentByIDAndUserID(ctx context.Context, id, userID uint64) (*domain.Comment, error) {
	var result *domain.Comment

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find comment by id and user id: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) CreateComment(ctx context.Context, data *domain.Comment) (*domain.Comment, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.Debug().WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on create comment: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) UpdateComment(ctx context.Context, data *domain.Comment) (*domain.Comment, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Updates(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on update comment: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) DeleteCommentByID(ctx context.Context, id uint64) error {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Comment{}).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on delete comment by id : %v", err)

		return sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return nil
}
