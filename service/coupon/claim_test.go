package coupon

import (
	m "coupon_be/mock"
	"coupon_be/request"
	sharedErrs "coupon_be/shared/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (suite *CouponServiceTestSuite) Test_Claim() {
	user := m.InitUserDomain()
	coupon := m.InitCouponDomain()
	input := &request.ClaimCoupon{
		CouponName: "COUPON_TEST",
		Username:   "user_123",
	}

	testCases := []struct {
		name          string
		prepareMock   func()
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			prepareMock: func() {
				suite.repo.EXPECT().FindCouponByName(suite.ctx, gomock.Eq(input.CouponName), gomock.Eq(false)).
					Return(coupon, nil).
					Times(1)
				suite.repo.EXPECT().FindCountUserClaimByCouponID(suite.ctx, gomock.Eq(coupon.ID)).
					Return(int64(0), nil).
					Times(1)
				suite.repo.EXPECT().FindUserByUsername(suite.ctx, gomock.Eq(input.Username)).
					Return(user, nil).
					Times(1)
				suite.repo.EXPECT().FindUserClaimByUserIDAndCouponID(suite.ctx, gomock.Eq(coupon.ID), gomock.Eq(user.ID)).
					Return(nil, nil).
					Times(1)
				suite.repo.EXPECT().CreateUserClaim(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
				suite.repo.EXPECT().DecrementCouponRemainingAmount(gomock.Any(), gomock.Eq(coupon.ID)).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "coupon not found",
			prepareMock: func() {
				suite.repo.EXPECT().FindCouponByName(suite.ctx, gomock.Eq(input.CouponName), gomock.Eq(false)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().FindUserByUsername(suite.ctx, gomock.Eq(input.Username)).
					Times(0)
				suite.repo.EXPECT().FindUserClaimByUserIDAndCouponID(suite.ctx, gomock.Eq(coupon.ID), gomock.Eq(user.ID)).
					Times(0)
				suite.repo.EXPECT().CreateUserClaim(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().DecrementCouponRemainingAmount(gomock.Any(), gomock.Eq(coupon.ID)).
					Times(0)
			},
			wantErr:       true,
			expectedError: sharedErrs.NotFoundErr,
		},
		{
			name: "coupon not usable",
			prepareMock: func() {
				c := m.InitCouponDomain()
				c.RemainingAmount = 0

				suite.repo.EXPECT().FindCouponByName(suite.ctx, gomock.Eq(input.CouponName), gomock.Eq(false)).
					Return(c, nil).
					Times(1)
				suite.repo.EXPECT().FindCountUserClaimByCouponID(suite.ctx, gomock.Eq(coupon.ID)).
					Return(int64(c.Amount), nil).
					Times(1)
				suite.repo.EXPECT().FindUserByUsername(suite.ctx, gomock.Eq(input.Username)).
					Times(0)
				suite.repo.EXPECT().FindUserClaimByUserIDAndCouponID(suite.ctx, gomock.Eq(coupon.ID), gomock.Eq(user.ID)).
					Times(0)
				suite.repo.EXPECT().CreateUserClaim(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().DecrementCouponRemainingAmount(gomock.Any(), gomock.Eq(coupon.ID)).
					Times(0)
			},
			wantErr: true,
			expectedError: sharedErrs.NewBusinessValidationErr(
				"Coupon %s is not usable because no stock remaining", coupon.Name),
		},
		{
			name: "user not found",
			prepareMock: func() {
				suite.repo.EXPECT().FindCouponByName(suite.ctx, gomock.Eq(input.CouponName), gomock.Eq(false)).
					Return(coupon, nil).
					Times(1)
				suite.repo.EXPECT().FindCountUserClaimByCouponID(suite.ctx, gomock.Eq(coupon.ID)).
					Return(int64(0), nil).
					Times(1)
				suite.repo.EXPECT().FindUserByUsername(suite.ctx, gomock.Eq(input.Username)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().FindUserClaimByUserIDAndCouponID(suite.ctx, gomock.Eq(coupon.ID), gomock.Eq(user.ID)).
					Times(0)
				suite.repo.EXPECT().CreateUserClaim(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().DecrementCouponRemainingAmount(gomock.Any(), gomock.Eq(coupon.ID)).
					Times(0)
			},
			wantErr:       true,
			expectedError: sharedErrs.NotFoundErr,
		},
		{
			name: "user already claimed",
			prepareMock: func() {
				suite.repo.EXPECT().FindCouponByName(suite.ctx, gomock.Eq(input.CouponName), gomock.Eq(false)).
					Return(coupon, nil).
					Times(1)
				suite.repo.EXPECT().FindCountUserClaimByCouponID(suite.ctx, gomock.Eq(coupon.ID)).
					Return(int64(0), nil).
					Times(1)
				suite.repo.EXPECT().FindUserByUsername(suite.ctx, gomock.Eq(input.Username)).
					Return(user, nil).
					Times(1)
				suite.repo.EXPECT().FindUserClaimByUserIDAndCouponID(suite.ctx, gomock.Eq(coupon.ID), gomock.Eq(user.ID)).
					Return(m.InitUserClaimDomain(), nil).
					Times(1)
				suite.repo.EXPECT().CreateUserClaim(gomock.Any(), gomock.Any()).
					Times(0)
				suite.repo.EXPECT().DecrementCouponRemainingAmount(gomock.Any(), gomock.Eq(coupon.ID)).
					Times(0)
			},
			wantErr: true,
			expectedError: sharedErrs.New(sharedErrs.ErrKindConflict, "Coupon %s is already claimed by user %s",
				coupon.Name, user.Username),
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			defer suite.After(t)
			tc.prepareMock()

			// Act
			err := suite.couponService.Claim(suite.ctx, input)

			// Assert
			assert.Equal(t, tc.wantErr, err != nil, "error expected %v, but actual: %v", tc.wantErr, err)
			if tc.wantErr {
				assert.Error(t, err)
			}
		})
	}
}
