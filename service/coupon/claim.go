package coupon

import (
	"context"
	"coupon_be/domain"
	"coupon_be/request"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/external/database"
	"coupon_be/util"
	"coupon_be/util/logger"
	"database/sql"
	"errors"
	"time"
)

func (b *base) Claim(ctx context.Context, input *request.ClaimCoupon) error {
	logger.Info(ctx, "Claim Coupon with req: %v", input)

	coupon, err := b.repository.FindCouponByName(ctx, util.SanitizeString(input.CouponName), false)
	if err != nil {
		return err
	}

	logger.Info(ctx, "resync coupon %s remaining amount ...", coupon.Name)
	if err = b.resyncCouponRemainingAmount(ctx, coupon); err != nil {
		return err
	}

	logger.Info(ctx, "coupon %s usable amount: %d remaining", coupon.Name, coupon.RemainingAmount)
	if valid := coupon.IsUsable(); !valid {
		logger.Warn(ctx, "coupon %s is not usable", coupon.Name)
		return sharedErrs.NewBusinessValidationErr("Coupon %s is not usable because no stock remaining", coupon.Name)
	}

	user, err := b.repository.FindUserByUsername(ctx, input.Username)
	if err != nil {
		return err
	}

	logger.Debug(ctx, "checking whether coupon %s is already claimed by user id %d ...", input.CouponName, user.ID)
	userClaimExists, err := b.repository.FindUserClaimByUserIDAndCouponID(ctx, user.ID, coupon.ID)
	if err != nil && !errors.Is(err, sharedErrs.NotFoundErr) {
		return err
	}
	if userClaimExists != nil {
		logger.Warn(ctx, "coupon %s is already claimed by user id %d with user claim id %d", coupon.Name, user.ID, userClaimExists.ID)

		return sharedErrs.New(sharedErrs.ErrKindConflict, "Coupon %s is already claimed by user %s",
			coupon.Name, user.Username)
	}

	logger.Info(ctx, "claiming coupon %s for user id %d ...", coupon.Name, user.ID)

	tCtx, tx := database.InitTx(ctx, b.writeDB)
	defer func() {
		if err = tx.Rollback().Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(ctx, "Repository Error on executing b.Claim: ROLLBACK TXN: %v", err)
		}
	}()

	now := time.Now()
	_, err = b.repository.CreateUserClaim(tCtx, &domain.UserClaim{
		BaseModel: domain.BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:   user.ID,
		CouponID: coupon.ID,
	})
	if err != nil {
		return err
	}

	if err = b.repository.DecrementCouponRemainingAmount(tCtx, coupon.ID); err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		logger.Error(ctx, "Repository Error on executing b.Claim: COMMIT TXN: %v", err)
	}

	return nil
}

func (b *base) resyncCouponRemainingAmount(ctx context.Context, coupon *domain.Coupon) error {
	claimCount, err := b.repository.FindUserClaimCountByCouponID(ctx, coupon.ID)
	if err != nil {
		return err
	}

	actualRemainingAmount := coupon.Amount - uint64(claimCount)

	logger.Info(ctx, "coupon %d amount: %d | remaining amount: %d | claim count: %d | actual remaining: %d",
		coupon.ID, claimCount, coupon.Amount, coupon.RemainingAmount, actualRemainingAmount)
	if actualRemainingAmount != coupon.RemainingAmount {
		logger.Info(ctx, "updating coupon %d remaining amount to %d", coupon.ID, actualRemainingAmount)

		coupon.RemainingAmount = actualRemainingAmount
		if _, err = b.repository.UpdateCoupon(ctx, coupon); err != nil {
			return err
		}
	} else {
		logger.Debug(ctx, "coupon %d remaining amount is correct", coupon.ID)
	}

	return nil
}
