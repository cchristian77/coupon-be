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

func (b base) Claim(ctx context.Context, input *request.ClaimCoupon) error {
	logger.Info(ctx, "Claim Coupon with req: %v", input)

	tCtx, tx := database.InitTx(ctx, b.writeDB)
	defer func() {
		if err := tx.Rollback().Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(ctx, "Repository Error on executing b.Claim: ROLLBACK TXN: %v", err)
		}
	}()

	coupon, err := b.repository.FindCouponByName(tCtx, util.SanitizeString(input.CouponName), false)
	if err != nil {
		return err
	}

	user, err := b.repository.FindUserByUsername(tCtx, input.UserName)
	if err != nil {
		return err
	}

	logger.Info(ctx, "coupon %s usable amount: %d", coupon.Name, coupon.RemainingAmount)
	if isUsable := coupon.IsUsable(); !isUsable {
		logger.Error(ctx, "coupon %s is not usable", coupon.Name)
		return sharedErrs.NewBusinessValidationErr("%s", "Coupon COUPON_TEST is not usable because no stock remaining")
	}

	logger.Debug(ctx, "checking whether coupon %s is already claimed by user id %d ...", input.CouponName, user.ID)
	userClaimExists, err := b.repository.FindUserClaimByUserIDAndCouponID(tCtx, user.ID, coupon.ID)
	if err != nil && !errors.Is(err, sharedErrs.NotFoundErr) {
		return err
	}
	if userClaimExists != nil {
		logger.Warn(ctx, "coupon %s is already claimed by user id %d with user claim id %d", coupon.Name, user.ID, userClaimExists.ID)

		return sharedErrs.New(sharedErrs.ErrKindConflict, "Coupon %s is already claimed by user %s",
			coupon.Name, user.Username)
	}

	logger.Debug(ctx, "claiming coupon %s for user id %d ...", coupon.Name, user.ID)
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

	logger.Debug(ctx, "decrementing coupon %s remaining amount to %d ...", coupon.Name, coupon.RemainingAmount-1)
	if _, err = b.repository.DecrementCouponRemainingAmount(tCtx, coupon.ID); err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		logger.Error(ctx, "Repository Error on executing b.Claim: COMMIT TXN: %v", err)
	}

	return nil
}
