package comment

import (
	"coupon_be/domain"
	m "coupon_be/mock"
	"coupon_be/response"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (suite *CommentServiceTestSuite) Test_Detail() {
	var (
		comment   *domain.Comment
		expected  *response.Comment
		commentID uint64 = 1
	)

	testCases := []struct {
		name          string
		prepareMock   func()
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			prepareMock: func() {
				comment = m.InitCommentDomain()
				expected = response.NewCommentFromDomain(comment)

				suite.repo.EXPECT().FindCommentByID(suite.ctx, gomock.Eq(commentID)).
					Return(comment, nil).
					Times(1)
			},
		},
		{
			name: "data not found",
			prepareMock: func() {
				suite.repo.EXPECT().FindCommentByID(suite.ctx, gomock.Eq(commentID)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
			},
			wantErr:       true,
			expectedError: sharedErrs.NotFoundErr,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			defer suite.After(t)
			tc.prepareMock()

			// Act
			result, err := suite.commentService.Detail(suite.ctx, commentID)

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
