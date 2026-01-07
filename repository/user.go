package repository

import (
	"context"
	"coupon_be/domain"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/external/database"
	"coupon_be/util/logger"

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
