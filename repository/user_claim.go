package repository

import (
	"context"
	"coupon_be/domain"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/external/database"
	"coupon_be/util/logger"

	"gorm.io/gorm/clause"
)

func (r *repo) FindUserClaimByUserIDAndCouponID(ctx context.Context, userID, couponID uint64) (*domain.UserClaim, error) {
	var result *domain.UserClaim

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("user_id = ? AND coupon_id = ?", userID, couponID).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find user claim by user id and coupon id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) CreateUserClaim(ctx context.Context, data *domain.UserClaim) (*domain.UserClaim, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.Debug().WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on create user claim: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}
