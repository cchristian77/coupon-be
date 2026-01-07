package comment

import (
	m "coupon_be/mock"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/util/constant"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (suite *CommentServiceTestSuite) Test_Delete() {
	comment := m.InitCommentDomain()

	testCases := []struct {
		name          string
		prepareMock   func()
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			prepareMock: func() {
				authUser := constant.AuthUserFromCtx(suite.ctx)

				suite.repo.EXPECT().FindCommentByIDAndUserID(suite.ctx, gomock.Eq(comment.ID), gomock.Eq(authUser.ID)).
					Return(comment, nil).
					Times(1)
				suite.repo.EXPECT().DeleteCommentByID(suite.ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "post not found",
			prepareMock: func() {
				authUser := constant.AuthUserFromCtx(suite.ctx)

				suite.repo.EXPECT().FindCommentByIDAndUserID(suite.ctx, gomock.Eq(comment.ID), gomock.Eq(authUser.ID)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().DeleteCommentByID(suite.ctx, gomock.Eq(comment.ID)).
					Times(0)
			},
			wantErr:       true,
			expectedError: sharedErrs.NotFoundErr,
		},
		{
			name: "unexpected error",
			prepareMock: func() {
				authUser := constant.AuthUserFromCtx(suite.ctx)

				suite.repo.EXPECT().FindCommentByIDAndUserID(suite.ctx, gomock.Eq(comment.ID), gomock.Eq(authUser.ID)).
					Return(comment, nil).
					Times(1)
				suite.repo.EXPECT().DeleteCommentByID(suite.ctx, gomock.Eq(comment.ID)).
					Return(sharedErrs.InternalServerErr).
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
			err := suite.commentService.Delete(suite.ctx, comment.ID)

			// Assert
			assert.Equal(t, tc.wantErr, err != nil, "error expected %v, but actual: %v", tc.wantErr, err)
			if tc.wantErr {
				assert.Error(t, err)
			}
		})
	}
}
