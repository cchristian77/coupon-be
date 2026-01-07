package comment

import (
	m "base_project/mock"
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestNewService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := m.NewMockRepository(ctrl)
	writeDB, _, err := m.NewMockDB()
	if err != nil {
		t.Fatal(err)
	}

	attendanceService, err := NewService(repoMock, writeDB)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, attendanceService)
	assert.Implements(t, (*Service)(nil), attendanceService)
}

type CommentServiceTestSuite struct {
	suite.Suite
	repo    *m.MockRepository
	writeDB *gorm.DB
	sqlMock sqlmock.Sqlmock
	ctx     context.Context

	commentService Service
}

func (suite *CommentServiceTestSuite) Before(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var err error

	logger.Initialise()

	suite.ctx = context.Background()
	suite.ctx = context.WithValue(suite.ctx, constant.XAuthUserKey, m.InitUserDomain())
	suite.repo = m.NewMockRepository(ctrl)
	suite.writeDB, suite.sqlMock, err = m.NewMockDB()
	if err != nil {
		t.Fatal(err)
	}

	suite.commentService, err = NewService(suite.repo, suite.writeDB)
	if err != nil {
		t.Fatal(err)
	}
}

func (suite *CommentServiceTestSuite) After(t *testing.T) {}

func TestSuiteRunCommentService(t *testing.T) {
	suite.Run(t, new(CommentServiceTestSuite))
}
