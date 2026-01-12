package coupon

import (
	m "coupon_be/mock"
	"coupon_be/request"
	"coupon_be/response"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (suite *CouponServiceTestSuite) Test_Store() {
	coupon := m.InitCouponDomain()
	input := &request.UpsertCoupon{
		ID:     1,
		Name:   "COUPON_TEST",
		Amount: 50,
	}

	var expected *response.Coupon

	testCases := []struct {
		name          string
		prepareMock   func()
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			prepareMock: func() {
				expected = response.NewCouponFromDomain(coupon)

				suite.repo.EXPECT().FindCouponByName(suite.ctx, gomock.Eq(input.Name), gomock.Eq(false)).
					Return(nil, nil).
					Times(1)
				suite.repo.EXPECT().CreateCoupon(suite.ctx, gomock.Any()).
					Return(coupon, nil).
					Times(1)
			},
		},
		{
			name: "coupon name exists",
			prepareMock: func() {
				suite.repo.EXPECT().FindCouponByName(suite.ctx, gomock.Eq(input.Name), gomock.Eq(false)).
					Return(coupon, nil).
					Times(1)
				suite.repo.EXPECT().CreateCoupon(suite.ctx, gomock.Any()).
					Times(0)
			},
			wantErr: true,
			expectedError: sharedErrs.NewBusinessValidationErr(
				"Create Failed. Coupon with name '%s' already exists.", input.Name),
		},
		{
			name: "unexpected error",
			prepareMock: func() {
				suite.repo.EXPECT().FindCouponByName(suite.ctx, gomock.Eq(input.Name), gomock.Eq(false)).
					Return(coupon, nil).
					Times(1)
				suite.repo.EXPECT().CreateCoupon(suite.ctx, gomock.Any()).
					Return(nil, sharedErrs.InternalServerErr).
					Times(1)
			},
			wantErr:       true,
			expectedError: sharedErrs.InternalServerErr,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			defer suite.After(t)
			tc.prepareMock()

			// Act
			result, err := suite.couponService.Store(suite.ctx, input)

			// Assert
			assert.Equal(t, tc.wantErr, err != nil, "error expected %v, but actual: %v", tc.wantErr, err)
			if tc.wantErr {
				assert.Empty(t, result)
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, result)
				if err = util.CompareData(result, expected, 1); err != nil {
					t.Errorf("error on comparing data : %v", err)
				}
			}
		})
	}
}
