package comment

import (
	m "base_project/mock"
	"base_project/request"
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util"
	"base_project/util/constant"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (suite *CommentServiceTestSuite) Test_Update() {
	comment := m.InitCommentDomain()
	input := &request.UpdateComment{
		ID:      1,
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
				authUser := constant.AuthUserFromCtx(suite.ctx)
				expected = response.NewCommentFromDomain(comment)

				suite.repo.EXPECT().FindCommentByIDAndUserID(suite.ctx, gomock.Eq(input.ID), gomock.Eq(authUser.ID)).
					Return(comment, nil).
					Times(1)
				suite.repo.EXPECT().UpdateComment(suite.ctx, gomock.Any()).
					Return(comment, nil).
					Times(1)
			},
		},
		{
			name: "post not found",
			prepareMock: func() {
				authUser := constant.AuthUserFromCtx(suite.ctx)

				suite.repo.EXPECT().FindCommentByIDAndUserID(suite.ctx, gomock.Eq(input.ID), gomock.Eq(authUser.ID)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().UpdateComment(suite.ctx, gomock.Any()).
					Times(0)
			},
			wantErr:       true,
			expectedError: sharedErrs.NotFoundErr,
		},
		{
			name: "unexpected error",
			prepareMock: func() {
				authUser := constant.AuthUserFromCtx(suite.ctx)

				suite.repo.EXPECT().FindCommentByIDAndUserID(suite.ctx, gomock.Eq(input.ID), gomock.Eq(authUser.ID)).
					Return(comment, nil).
					Times(1)
				suite.repo.EXPECT().UpdateComment(suite.ctx, gomock.Any()).
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
			result, err := suite.commentService.Update(suite.ctx, input)

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
