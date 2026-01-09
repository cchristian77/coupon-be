package repository

import (
	"context"
	"coupon_be/domain"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/external/database"
	"coupon_be/util/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *repo) FindCouponByID(ctx context.Context, id uint64) (*domain.Coupon, error) {
	var result *domain.Coupon

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		First(&result, id).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find coupon by id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindCouponByName(ctx context.Context, name string, withClaimBy bool) (*domain.Coupon, error) {
	var result *domain.Coupon

	db, _ := database.ConnFromCtx(ctx, r.DB)

	query := db.WithContext(ctx)
	if withClaimBy {
		query.Preload("ClaimBy")
	}

	err := query.Where("name = ?", name).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find coupon by name : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindCoupons(ctx context.Context, search string) ([]*domain.Coupon, error) {
	var result []*domain.Coupon

	db, _ := database.ConnFromCtx(ctx, r.DB)

	query := db.WithContext(ctx).Model(&result)

	if search != "" {
		query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&result).Error; err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find coupons : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) CreateCoupon(ctx context.Context, data *domain.Coupon) (*domain.Coupon, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on create coupon: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) UpdateCoupon(ctx context.Context, data *domain.Coupon) (*domain.Coupon, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Updates(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on update coupon: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) DecrementCouponRemainingAmount(ctx context.Context, id uint64) (*domain.Coupon, error) {
	var result *domain.Coupon

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Model(&result).
		Where("id = ? AND remaining_amount > 0", id).
		UpdateColumn("remaining_amount", gorm.Expr("remaining_amount - 1")).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on update coupon: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}
