package comment

import (
	"bytes"
	"context"
	m "coupon_be/mock"
	ms "coupon_be/mock/service"
	"coupon_be/request"
	"coupon_be/response"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/fhttp"
	"coupon_be/shared/fhttp/middleware"
	"coupon_be/util"
	"coupon_be/util/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CommentControllerTestSuite struct {
	suite.Suite
	router         *mux.Router
	commentService *ms.MockCommentService
	repo           *m.MockRepository
	controller     *Controller
	prefix         string
	ctx            context.Context
}

func (suite *CommentControllerTestSuite) Before(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite.router = mux.NewRouter()
	suite.router.Use(middleware.PanicRecovery())
	suite.commentService = ms.NewMockCommentService(ctrl)
	suite.repo = m.NewMockRepository(ctrl)
	suite.controller = &Controller{
		comment: suite.commentService,
	}

	suite.controller.RegisterRoutes(suite.router.PathPrefix(suite.prefix).Subrouter())
}

func (suite *CommentControllerTestSuite) Test_FilterComments() {
	filterCommentReq := &request.FilterComment{
		PostID:  1,
		Page:    1,
		PerPage: 10,
		Search:  "search",
	}

	var comments []*response.Comment
	for i := 0; i < 5; i++ {
		c := m.InitCommentResponse()
		c.ID = uint64(i + 1)
		comments = append(comments, c)
	}
	expected := comments

	type testCase struct {
		name           string
		httpMethod     string
		input          *request.FilterComment
		requestHeaders struct {
			authorization string
		}
		prepareMock        func(t *testing.T)
		expectedStatusCode int
	}
	testCases := []testCase{
		{
			name:       "success",
			httpMethod: http.MethodGet,
			input:      filterCommentReq,
			prepareMock: func(t *testing.T) {
				p := &util.Pagination{}
				p.SetPage(filterCommentReq.Page)
				p.SetPage(filterCommentReq.PerPage)

				paginatedComments := response.NewBasePagination(comments, p)

				suite.commentService.EXPECT().
					FilterComments(gomock.Any(), gomock.Any()).
					Return(paginatedComments, nil).
					Times(1)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "unauthorized",
			httpMethod: http.MethodGet,
			input:      filterCommentReq,
			requestHeaders: struct{ authorization string }{
				authorization: "",
			},
			prepareMock: func(t *testing.T) {
				suite.commentService.EXPECT().
					FilterComments(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:       "unexpected error",
			httpMethod: http.MethodGet,
			input:      filterCommentReq,
			prepareMock: func(t *testing.T) {
				suite.commentService.EXPECT().
					FilterComments(gomock.Any(), gomock.Any()).
					Return(nil, sharedErrs.InternalServerErr).
					Times(1)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.prepareMock(t)
			r := httptest.NewRecorder()
			url := fmt.Sprintf("%s/%d/posts?page=%v&per_page=%v&search=%v",
				suite.prefix, tc.input.PostID, tc.input.Page, tc.input.PerPage, tc.input.Search)

			// Act
			req, err := http.NewRequestWithContext(suite.ctx, tc.httpMethod, url, nil)
			req.Header.Add("Authorization", "Bearer "+tc.requestHeaders.authorization)
			if err != nil {
				t.Fatalf("unit test for %s endpoint failed: %v", url, err)
			}
			suite.router.ServeHTTP(r, req)

			// Assert
			resp := r.Result()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedStatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("error reading response body: %v", err)
				}

				httpResp := &fhttp.Response{Data: &response.BasePagination[[]*response.Comment]{}}
				if err = json.Unmarshal(body, httpResp); err != nil {
					t.Errorf("error to unmarshal response body: %v", err)
				}

				result := httpResp.Data.(*response.BasePagination[[]*response.Comment])
				for i, actual := range result.Data {
					if err = util.CompareData(expected[i], actual, 1); err != nil {
						t.Errorf("error on comparing data : %v", err)
					}
				}
			}
		})
	}
}

func (suite *CommentControllerTestSuite) Test_Detail() {
	commentID := uint64(1)
	commentResponse := m.InitCommentResponse()

	var expectedResult *response.Comment

	type testCase struct {
		name       string
		httpMethod string
		input      struct {
			commentID any
		}
		requestHeaders struct {
			authorization string
		}
		prepareMock        func(t *testing.T)
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			name:       "success",
			httpMethod: http.MethodGet,
			input:      struct{ commentID any }{commentID: commentID},
			prepareMock: func(t *testing.T) {
				expectedResult = commentResponse
				suite.commentService.EXPECT().
					Detail(gomock.Any(), gomock.Eq(commentID)).
					Return(commentResponse, nil).
					Times(1)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "invalid comment id",
			httpMethod: http.MethodGet,
			input:      struct{ commentID any }{commentID: "invalid"},
			prepareMock: func(t *testing.T) {
				suite.commentService.EXPECT().
					Detail(gomock.Any(), gomock.Eq(commentID)).
					Times(0)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:       "unauthorized",
			httpMethod: http.MethodGet,
			input:      struct{ commentID any }{commentID: commentID},
			requestHeaders: struct{ authorization string }{
				authorization: "",
			},
			prepareMock: func(t *testing.T) {
				suite.commentService.EXPECT().
					Detail(gomock.Any(), gomock.Eq(commentID)).
					Times(0)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:       "method not allowed",
			httpMethod: http.MethodPatch,
			input:      struct{ commentID any }{commentID: commentID},
			requestHeaders: struct{ authorization string }{
				authorization: "",
			},
			prepareMock: func(t *testing.T) {
				suite.commentService.EXPECT().
					Detail(gomock.Any(), gomock.Eq(commentID)).
					Times(0)
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "not found error",
			httpMethod: http.MethodGet,
			input:      struct{ commentID any }{commentID: commentID},
			prepareMock: func(t *testing.T) {
				expectedResult = commentResponse
				suite.commentService.EXPECT().
					Detail(gomock.Any(), gomock.Eq(commentID)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:       "unexpected error",
			httpMethod: http.MethodGet,
			input:      struct{ commentID any }{commentID: commentID},
			prepareMock: func(t *testing.T) {
				expectedResult = commentResponse
				suite.commentService.EXPECT().
					Detail(gomock.Any(), gomock.Eq(commentID)).
					Return(nil, sharedErrs.InternalServerErr).
					Times(1)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.prepareMock(t)
			r := httptest.NewRecorder()
			url := fmt.Sprintf("%s/%v", suite.prefix, tc.input.commentID)

			// Act
			req, err := http.NewRequestWithContext(suite.ctx, tc.httpMethod, url, nil)
			req.Header.Add("Authorization", "Bearer "+tc.requestHeaders.authorization)
			if err != nil {
				t.Fatalf("unit test for %s endpoint failed: %v", url, err)
			}
			suite.router.ServeHTTP(r, req)

			// Assert
			resp := r.Result()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedStatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("error reading response body: %v", err)
				}

				httpResp := &fhttp.Response{Data: &response.Comment{}}
				if err = json.Unmarshal(body, httpResp); err != nil {
					t.Errorf("error to unmarshal response body: %v", err)
				}

				result := httpResp.Data.(*response.Comment)
				if err = util.CompareData(expectedResult, result, 1); err != nil {
					t.Errorf("error on comparing data : %v", err)
				}
			}
		})
	}
}

func (suite *CommentControllerTestSuite) Test_Store() {
	createCommentReq := &request.CreateComment{
		PostID:  1,
		Comment: "comment test",
		Rating:  util.ToPointerValue(uint8(5)),
	}
	comment := m.InitCommentDomain()

	var expectedResult *response.Comment

	type testCase struct {
		name           string
		httpMethod     string
		input          any
		requestHeaders struct {
			authorization string
		}
		prepareMock        func(t *testing.T)
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			name:       "success",
			httpMethod: http.MethodPost,
			input:      createCommentReq,
			prepareMock: func(t *testing.T) {
				expectedResult = response.NewCommentFromDomain(comment)

				suite.commentService.EXPECT().
					Store(gomock.Any(), gomock.Any()).
					Return(expectedResult, nil).
					Times(1)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "validation error",
			httpMethod: http.MethodPost,
			input:      nil,
			prepareMock: func(t *testing.T) {
				expectedResult = response.NewCommentFromDomain(comment)

				suite.commentService.EXPECT().
					Store(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.prepareMock(t)
			r := httptest.NewRecorder()
			url := fmt.Sprintf("%s", suite.prefix)

			reqBody := new(bytes.Buffer)
			if err := json.NewEncoder(reqBody).Encode(tc.input); err != nil {
				t.Fatalf("error to encode request to bytes: %v", err)
			}

			// Act
			req, err := http.NewRequestWithContext(suite.ctx, tc.httpMethod, url, reqBody)
			req.Header.Add("Authorization", "Bearer "+tc.requestHeaders.authorization)
			if err != nil {
				t.Fatalf("unit test for %s endpoint failed: %v", url, err)
			}
			suite.router.ServeHTTP(r, req)

			// Assert
			resp := r.Result()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedStatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("error reading response body: %v", err)
				}

				httpResp := &fhttp.Response{Data: &response.Comment{}}
				if err = json.Unmarshal(body, &httpResp); err != nil {
					t.Errorf("error to unmarshal response body: %v", err)
				}

				result := httpResp.Data.(*response.Comment)
				if err = util.CompareData(expectedResult, result, 1); err != nil {
					t.Errorf("error on comparing data : %v", err)
				}
			}
		})
	}
}

func (suite *CommentControllerTestSuite) Test_Update() {
	updateCommentReq := &request.UpdateComment{
		ID:      1,
		Comment: "comment test",
		Rating:  util.ToPointerValue(uint8(5)),
	}
	comment := m.InitCommentDomain()

	var expectedResult *response.Comment

	type testCase struct {
		name       string
		httpMethod string
		input      struct {
			commentID any
			req       any
		}
		requestHeaders struct {
			authorization string
		}
		prepareMock        func(t *testing.T)
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			name:       "success",
			httpMethod: http.MethodPut,
			input: struct {
				commentID any
				req       any
			}{commentID: comment.ID, req: updateCommentReq},
			prepareMock: func(t *testing.T) {
				expectedResult = response.NewCommentFromDomain(comment)

				suite.commentService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(expectedResult, nil).
					Times(1)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "invalid comment id",
			httpMethod: http.MethodPut,
			input: struct {
				commentID any
				req       any
			}{commentID: "invalid", req: updateCommentReq},
			prepareMock: func(t *testing.T) {
				suite.commentService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.prepareMock(t)
			r := httptest.NewRecorder()
			url := fmt.Sprintf("%s/%v", suite.prefix, tc.input.commentID)

			reqBody := new(bytes.Buffer)
			if err := json.NewEncoder(reqBody).Encode(tc.input.req); err != nil {
				t.Fatalf("error to encode request to bytes: %v", err)
			}

			// Act
			req, err := http.NewRequestWithContext(suite.ctx, tc.httpMethod, url, reqBody)
			req.Header.Add("Authorization", "Bearer "+tc.requestHeaders.authorization)
			if err != nil {
				t.Fatalf("unit test for %s endpoint failed: %v", url, err)
			}
			suite.router.ServeHTTP(r, req)

			// Assert
			resp := r.Result()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedStatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("error reading response body: %v", err)
				}

				httpResp := &fhttp.Response{Data: &response.Comment{}}
				if err = json.Unmarshal(body, &httpResp); err != nil {
					t.Errorf("error to unmarshal response body: %v", err)
				}

				result := httpResp.Data.(*response.Comment)
				if err = util.CompareData(expectedResult, result, 1); err != nil {
					t.Errorf("error on comparing data : %v", err)
				}
			}
		})
	}
}

func (suite *CommentControllerTestSuite) Test_Delete() {
	commentID := uint64(1)

	type testCase struct {
		name           string
		httpMethod     string
		input          any
		requestHeaders struct {
			authorization string
		}
		prepareMock        func(t *testing.T)
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			name:       "success",
			httpMethod: http.MethodDelete,
			input:      commentID,
			prepareMock: func(t *testing.T) {

				suite.commentService.EXPECT().
					Delete(gomock.Any(), gomock.Eq(commentID)).
					Return(nil).
					Times(1)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "invalid comment id",
			httpMethod: http.MethodDelete,
			input:      "invalid",
			prepareMock: func(t *testing.T) {
				suite.commentService.EXPECT().
					Delete(gomock.Any(), gomock.Eq(commentID)).
					Times(0)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:       "unexpect error",
			httpMethod: http.MethodDelete,
			input:      commentID,
			prepareMock: func(t *testing.T) {
				suite.commentService.EXPECT().
					Delete(gomock.Any(), gomock.Eq(commentID)).
					Return(sharedErrs.InternalServerErr).
					Times(1)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.prepareMock(t)
			r := httptest.NewRecorder()
			url := fmt.Sprintf("%s/%v", suite.prefix, tc.input)

			// Act
			req, err := http.NewRequestWithContext(suite.ctx, tc.httpMethod, url, nil)
			req.Header.Add("Authorization", "Bearer "+tc.requestHeaders.authorization)
			if err != nil {
				t.Fatalf("unit test for %s endpoint failed: %v", url, err)
			}
			suite.router.ServeHTTP(r, req)

			// Assert
			resp := r.Result()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestSuiteRunCommentController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger.Initialise()

	testSuite := new(CommentControllerTestSuite)
	testSuite.ctx = context.Background()
	testSuite.prefix = "/posts/v1"

	suite.Run(t, testSuite)
}
