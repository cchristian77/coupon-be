package repository

import (
	"coupon_be/domain"
	m "coupon_be/mock"
	"coupon_be/util"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	allCommentColNames = []string{
		`id`,
		`created_at`,
		`updated_at`,
		`post_id`,
		`user_id`,
		`comment`,
		`rating`,
	}
	allUserColNames = []string{
		`id`,
		`created_at`,
		`updated_at`,
		`deleted_at`,
		`username`,
		`password`,
		`email`,
		`full_name`,
		`role`,
	}
)

func (suite *RepositoryTestSuite) Test_FindCommentsPaginatedByPostID() {
	var (
		postID   uint64 = 2
		userID   uint64 = 1
		p        *util.Pagination
		expected []*domain.Comment
		userRow  = m.InitUserDomain()
	)
	p = new(util.Pagination)
	p.SetPage(1)
	p.SetLimit(10)

	var rows []driver.Value
	for i := 0; i < 1; i++ {
		c := m.InitCommentDomain()
		c.ID = uint64(i + 1)
		rows = append(rows, c.ID, c.CreatedAt, c.UpdatedAt, c.UserID, c.PostID, c.Comment, c.Rating)
		expected = append(expected, c)
	}

	type testCase struct {
		name        string
		mockClosure func(m sqlmock.Sqlmock)
		wantErr     bool
	}

	testCases := []testCase{
		{
			name: "success",
			mockClosure: func(m sqlmock.Sqlmock) {
				countStmt := m.ExpectQuery(`SELECT count`)
				countStmt.WithArgs(postID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(len(rows)))

				stmt := m.ExpectQuery(`SELECT (.+) FROM "comments"`)
				stmt.WithArgs(postID, p.Limit()).
					WillReturnRows(
						sqlmock.NewRows(allCommentColNames).AddRows(rows))

				preloadStmt := m.ExpectQuery(`SELECT (.+) FROM "users"`)
				preloadStmt.WithArgs(userID).
					WillReturnRows(
						sqlmock.NewRows(allUserColNames).AddRow(
							userRow.ID,
							userRow.CreatedAt,
							userRow.UpdatedAt,
							userRow.DeletedAt,
							userRow.Username,
							userRow.Password,
							userRow.Email,
							userRow.FullName,
							userRow.Role))
			},
		},
		{
			name: "unexpected error",
			mockClosure: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").
					WillReturnError(errors.New("unexpected error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.mockClosure(suite.sqlMock)

			// Act
			result, err := suite.repo.FindCommentsPaginatedByPostID(suite.ctx, postID, p)

			// Assert
			assert.Equal(t, tc.wantErr, err != nil, "error expected %v, but actual: %v", tc.wantErr, err)

			if tc.wantErr {
				assert.Empty(t, result)
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, result)
				for i, actual := range result {
					if err = util.CompareData(result[i], actual, 1); err != nil {
						t.Errorf("error on comparing data : %v", err)
					}
				}
			}

			if err = suite.sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func (suite *RepositoryTestSuite) Test_FindCommentsPaginatedByUserID() {
	var (
		userID   uint64 = 1
		p        *util.Pagination
		expected []*domain.Comment
		userRow  = m.InitUserDomain()
	)
	p = new(util.Pagination)
	p.SetPage(1)
	p.SetLimit(10)

	var rows []driver.Value
	for i := 0; i < 1; i++ {
		c := m.InitCommentDomain()
		c.ID = uint64(i + 1)
		rows = append(rows, c.ID, c.CreatedAt, c.UpdatedAt, c.UserID, c.PostID, c.Comment, c.Rating)
		expected = append(expected, c)
	}

	type testCase struct {
		name        string
		mockClosure func(m sqlmock.Sqlmock)
		wantErr     bool
	}

	testCases := []testCase{
		{
			name: "success",
			mockClosure: func(m sqlmock.Sqlmock) {
				countStmt := m.ExpectQuery(`SELECT count`)
				countStmt.WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(len(rows)))

				stmt := m.ExpectQuery(`SELECT (.+) FROM "comments"`)
				stmt.WithArgs(userID, p.Limit()).
					WillReturnRows(
						sqlmock.NewRows(allCommentColNames).AddRows(rows))

				preloadStmt := m.ExpectQuery(`SELECT (.+) FROM "users"`)
				preloadStmt.WithArgs(userID).
					WillReturnRows(
						sqlmock.NewRows(allUserColNames).AddRow(
							userRow.ID,
							userRow.CreatedAt,
							userRow.UpdatedAt,
							userRow.DeletedAt,
							userRow.Username,
							userRow.Password,
							userRow.Email,
							userRow.FullName,
							userRow.Role))
			},
		},
		{
			name: "unexpected error",
			mockClosure: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").
					WillReturnError(errors.New("unexpected error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.mockClosure(suite.sqlMock)

			// Act
			result, err := suite.repo.FindCommentsPaginatedByUserID(suite.ctx, userID, p)

			// Assert
			assert.Equal(t, tc.wantErr, err != nil, "error expected %v, but actual: %v", tc.wantErr, err)

			if tc.wantErr {
				assert.Empty(t, result)
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, result)
				for i, actual := range result {
					if err = util.CompareData(result[i], actual, 1); err != nil {
						t.Errorf("error on comparing data : %v", err)
					}
				}
			}

			if err = suite.sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func (suite *RepositoryTestSuite) Test_FindCommentByID() {
	var (
		commentID uint64 = 1
		row              = m.InitCommentDomain()
		limit            = 1
		expected  *domain.Comment
	)

	type testCase struct {
		name        string
		mockClosure func(m sqlmock.Sqlmock)
		wantErr     bool
	}

	testCases := []testCase{
		{
			name: "success",
			mockClosure: func(m sqlmock.Sqlmock) {
				expected = row

				m.ExpectQuery(`SELECT (.+) FROM "comments"`).
					WithArgs(commentID, limit).
					WillReturnRows(
						sqlmock.NewRows(allCommentColNames).AddRow(
							row.ID,
							row.CreatedAt,
							row.UpdatedAt,
							row.UserID,
							row.PostID,
							row.Comment,
							row.Rating))
			},
		},
		{
			name: "unexpected error",
			mockClosure: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT (.+) FROM "comments"`).
					WillReturnError(errors.New("unexpected error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.mockClosure(suite.sqlMock)

			// Act
			result, err := suite.repo.FindCommentByID(suite.ctx, commentID)

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

			if err = suite.sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func (suite *RepositoryTestSuite) Test_FindCommentByIDAndUserID() {
	var (
		commentID uint64 = 1
		userID    uint64 = 1
		row              = m.InitCommentDomain()
		limit            = 1
		expected  *domain.Comment
	)

	type testCase struct {
		name        string
		mockClosure func(m sqlmock.Sqlmock)
		wantErr     bool
	}

	testCases := []testCase{
		{
			name: "success",
			mockClosure: func(m sqlmock.Sqlmock) {
				expected = row

				m.ExpectQuery(`SELECT (.+) FROM "comments"`).
					WithArgs(commentID, userID, limit).
					WillReturnRows(
						sqlmock.NewRows(allCommentColNames).AddRow(
							row.ID,
							row.CreatedAt,
							row.UpdatedAt,
							row.UserID,
							row.PostID,
							row.Comment,
							row.Rating))
			},
		},
		{
			name: "unexpected error",
			mockClosure: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT (.+) FROM "comments"`).
					WillReturnError(errors.New("unexpected error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.mockClosure(suite.sqlMock)

			// Act
			result, err := suite.repo.FindCommentByIDAndUserID(suite.ctx, commentID, userID)

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

			if err = suite.sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func (suite *RepositoryTestSuite) Test_CreateComment() {
	var (
		row      = m.InitCommentDomain()
		data     = row
		expected *domain.Comment
	)

	type testCase struct {
		name        string
		mockClosure func(m sqlmock.Sqlmock)
		wantErr     bool
	}

	testCases := []testCase{
		{
			name: "success",
			mockClosure: func(m sqlmock.Sqlmock) {
				expected = row

				m.ExpectBegin()
				m.ExpectQuery(`INSERT INTO "comments"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), row.PostID, row.UserID, row.Comment, row.Rating, row.ID).
					WillReturnRows(
						sqlmock.NewRows(allCommentColNames).AddRow(
							row.ID,
							row.CreatedAt,
							row.UpdatedAt,
							row.UserID,
							row.PostID,
							row.Comment,
							row.Rating))
				m.ExpectCommit()
			},
		},
		{
			name: "unexpected error",
			mockClosure: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectQuery(`INSERT INTO "comments"`).
					WillReturnError(errors.New("unexpected error"))
				m.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.mockClosure(suite.sqlMock)

			// Act
			result, err := suite.repo.CreateComment(suite.ctx, data)

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

			if err = suite.sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func (suite *RepositoryTestSuite) Test_UpdateComment() {
	var (
		row      = m.InitCommentDomain()
		data     = row
		expected *domain.Comment
	)

	type testCase struct {
		name        string
		mockClosure func(m sqlmock.Sqlmock)
		wantErr     bool
	}

	testCases := []testCase{
		{
			name: "success",
			mockClosure: func(m sqlmock.Sqlmock) {
				expected = row

				m.ExpectBegin()
				m.ExpectQuery(`UPDATE "comments"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), row.PostID, row.UserID, row.Comment, row.Rating, row.ID).
					WillReturnRows(
						sqlmock.NewRows(allCommentColNames).AddRow(
							row.ID,
							row.CreatedAt,
							row.UpdatedAt,
							row.UserID,
							row.PostID,
							row.Comment,
							row.Rating))
				m.ExpectCommit()
			},
		},
		{
			name: "unexpected error",
			mockClosure: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectQuery(`UPDATE "comments"`).
					WillReturnError(errors.New("unexpected error"))
				m.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.mockClosure(suite.sqlMock)

			// Act
			result, err := suite.repo.UpdateComment(suite.ctx, data)

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

			if err = suite.sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func (suite *RepositoryTestSuite) Test_DeleteCommentByID() {
	commentID := uint64(1)

	type testCase struct {
		name        string
		mockClosure func(m sqlmock.Sqlmock)
		wantErr     bool
	}

	testCases := []testCase{
		{
			name: "success",
			mockClosure: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`DELETE FROM "comments"`).
					WithArgs(commentID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
		},
		{
			name: "unexpected error",
			mockClosure: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectExec(`DELETE FROM "comments"`).
					WillReturnError(errors.New("unexpected error"))
				m.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			tc.mockClosure(suite.sqlMock)

			// Act
			err := suite.repo.DeleteCommentByID(suite.ctx, commentID)

			// Assert
			assert.Equal(t, tc.wantErr, err != nil, "error expected %v, but actual: %v", tc.wantErr, err)
			if tc.wantErr {
				assert.Error(t, err)
			}

			if err = suite.sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}
