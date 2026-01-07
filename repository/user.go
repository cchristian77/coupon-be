package repository

import (
	"base_project/domain"
	sharedErrs "base_project/shared/errors"
	"base_project/shared/external/database"
	"base_project/util/logger"
	"context"

	"gorm.io/gorm/clause"
)

func (r *repo) FindUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var result *domain.User

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("username = ?", username).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find user by username : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var result *domain.User

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("email = ?", email).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find user by email : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	var result *domain.User

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		First(&result, id).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find user by id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindBatchUsers(ctx context.Context, batchSize int, lastID uint64) ([]*domain.User, error) {
	var result []*domain.User

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("id > ?", lastID).
		Limit(batchSize).
		Find(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find batch users : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) CreateUser(ctx context.Context, data *domain.User) (*domain.User, error) {
	err := r.DB.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on create user : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}
