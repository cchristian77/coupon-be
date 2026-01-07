package comment

import (
	"coupon_be/domain"
	m "coupon_be/mock"
	"coupon_be/request"
	"coupon_be/response"
	"coupon_be/shared/errors"
	"coupon_be/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (suite *CommentServiceTestSuite) Test_Filter() {
	post := m.InitPostDomain()

	input := &request.FilterComment{
		PostID:  1,
		Page:    1,
		PerPage: 20,
	}

	var expected []*response.Comment

	comments := []*domain.Comment{m.InitCommentDomain()}
	for i := 0; i < 5; i++ {
		p := m.InitCommentDomain()
		p.ID = uint64(i + 1)
		comments = append(comments, p)
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
				for _, c := range comments {
					expected = append(expected, response.NewCommentFromDomain(c))
				}

				suite.repo.EXPECT().FindPostByID(suite.ctx, gomock.Eq(input.PostID), gomock.Eq(false)).
					Return(post, nil).
					Times(1)
				suite.repo.EXPECT().FindCommentsPaginatedByPostID(suite.ctx, gomock.Eq(input.PostID), gomock.Any()).
					Return(comments, nil).
					Times(1)
			},
		},
		{
			name: "post not found",
			prepareMock: func() {
				for _, c := range comments {
					expected = append(expected, response.NewCommentFromDomain(c))
				}

				suite.repo.EXPECT().FindPostByID(suite.ctx, gomock.Eq(input.PostID), gomock.Eq(false)).
					Return(nil, errors.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().FindCommentsPaginatedByPostID(suite.ctx, gomock.Eq(input.PostID), gomock.Any()).
					Times(0)
			},
			wantErr:       true,
			expectedError: errors.NewBusinessValidationErr("Post %d not found.", input.PostID),
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			defer suite.After(t)
			tc.prepareMock()

			// Act
			result, err := suite.commentService.FilterComments(suite.ctx, input)

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
