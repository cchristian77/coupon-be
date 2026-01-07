package repository

import (
	"base_project/domain"
	"base_project/domain/enums"
	sharedErrs "base_project/shared/errors"
	"base_project/shared/external/database"
	"base_project/util"
	"base_project/util/logger"
	"context"

	"gorm.io/gorm/clause"
)

func (r *repo) FindPublishedPostsPaginated(ctx context.Context, search string, p *util.Pagination) ([]*domain.Post, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	var (
		data  []*domain.Post
		count int64
	)

	query := db.WithContext(ctx).Debug().Model(&data).
		Preload("User").
		Where("status = ?", enums.PUBLISHEDPostStatus.String())

	if search != "" {
		query.Joins("JOIN users ON users.id = posts.user_id")
		query.Where("slug ILIKE ?", "%"+search+"%").
			Or("title ILIKE ?", "%"+search+"%").
			Or("users.full_name ILIKE ?", "%"+search+"%")

	}

	if err := query.Count(&count).Error; err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find post paginated count: %v", err)
		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	p.SetTotal(count)

	err := query.Offset(p.Offset()).
		Limit(p.Limit()).
		Find(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find post paginated: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) FindPostsPaginatedByUserID(ctx context.Context, userID uint64, search string, p *util.Pagination) ([]*domain.Post, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	var (
		data  []*domain.Post
		count int64
	)

	query := db.WithContext(ctx).Model(&data).Preload("User")

	if search != "" {
		query.Where("slug ILIKE ?", "%"+search+"%").
			Where("title ILIKE ?", "%"+search+"%")

		query.Joins("JOIN users ON users.id = posts.user_id AND users.full_name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find post paginated by user id count: %v", err)
		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	p.SetTotal(count)

	err := query.Offset(p.Offset()).
		Where("user_id = ?", userID).
		Limit(p.Limit()).
		Find(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find payslip by user id paginated : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) FindPostByID(ctx context.Context, id uint64, withComments bool) (*domain.Post, error) {
	var result *domain.Post

	db, _ := database.ConnFromCtx(ctx, r.DB)

	query := db.WithContext(ctx).Model(&result)

	if withComments {
		query.Preload("Comments")
	}

	if err := query.First(&result, id).Error; err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find post by id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindPostByIDAndUserID(ctx context.Context, id, userID uint64) (*domain.Post, error) {
	var result *domain.Post

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find post by id and user id: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindPostBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	var result *domain.Post

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find post by slug: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) CreatePost(ctx context.Context, data *domain.Post) (*domain.Post, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on create post: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) UpdatePost(ctx context.Context, data *domain.Post) (*domain.Post, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Updates(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on update post: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) UpdatePostStatus(ctx context.Context, id uint64, status enums.PostStatus) (*domain.Post, error) {
	result := &domain.Post{}

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Model(&result).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Update("status", status.String()).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on update post status: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) DeletePostByID(ctx context.Context, id uint64) error {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Post{}).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on delete post by id : %v", err)

		return sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return nil
}
