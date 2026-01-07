package comment

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

func (suite *CommentServiceTestSuite) Test_Store() {
	post := m.InitPostDomain()
	comment := m.InitCommentDomain()
	input := &request.CreateComment{
		PostID:  1,
		Comment: "test comment",
		Rating:  util.ToPointerValue(uint8(5)),
	}

	var expected *response.Comment

	testCases := []struct {
		name          string
		prepareMock   func()
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			prepareMock: func() {
				expected = response.NewCommentFromDomain(comment)

				suite.repo.EXPECT().FindPostByID(suite.ctx, gomock.Eq(input.PostID), gomock.Eq(false)).
					Return(post, nil).
					Times(1)
				suite.repo.EXPECT().CreateComment(suite.ctx, gomock.Any()).
					Return(comment, nil).
					Times(1)
			},
		},
		{
			name: "post not found",
			prepareMock: func() {
				suite.repo.EXPECT().FindPostByID(suite.ctx, gomock.Eq(input.PostID), gomock.Eq(false)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().CreateComment(suite.ctx, gomock.Any()).
					Times(0)
			},
			wantErr:       true,
			expectedError: sharedErrs.NotFoundErr,
		},
		{
			name: "unexpected error",
			prepareMock: func() {
				suite.repo.EXPECT().FindPostByID(suite.ctx, gomock.Eq(input.PostID), gomock.Eq(false)).
					Return(post, nil).
					Times(1)
				suite.repo.EXPECT().CreateComment(suite.ctx, gomock.Any()).
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
			result, err := suite.commentService.Store(suite.ctx, input)

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
