package coupon

import (
	"coupon_be/domain"
	m "coupon_be/mock"
	"coupon_be/request"
	"coupon_be/response"
	"coupon_be/util"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (suite *CouponServiceTestSuite) Test_Filter() {
	var (
		expected []*response.CouponList
		coupons  []*domain.Coupon
	)

	input := &request.FilterCoupon{
		Page:    1,
		PerPage: 20,
		Search:  "search_test",
	}

	for i := 0; i < 5; i++ {
		p := m.InitCouponDomain()
		p.ID = uint64(i + 1)
		coupons = append(coupons, p)
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
				for _, c := range coupons {
					expected = append(expected, response.NewCouponListFromDomain(c))
				}

				suite.repo.EXPECT().FindCouponsPaginated(suite.ctx, gomock.Eq(input.Search), gomock.Any()).
					Return(coupons, nil).
					Times(1)
			},
		},
		{
			name: "unexpected error",
			prepareMock: func() {
				suite.repo.EXPECT().FindCouponsPaginated(suite.ctx, gomock.Eq(input.Search), gomock.Any()).
					Return(nil, errors.New("unexpected error")).
					Times(1)
			},
			wantErr:       true,
			expectedError: errors.New("unexpected error"),
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			defer suite.After(t)
			tc.prepareMock()

			// Act
			result, err := suite.couponService.Filter(suite.ctx, input)

			// Assert
			assert.Equal(t, tc.wantErr, err != nil, "error expected %v, but actual: %v", tc.wantErr, err)
			if tc.wantErr {
				assert.Empty(t, result)
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, result)
				for i, actual := range result.Data {
					if err = util.CompareData(actual, expected[i], 1); err != nil {
						t.Fatalf("error on comparing data : %v", err)
					}
				}
			}
		})
	}
}
