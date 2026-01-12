package coupon

import (
	"context"
	m "coupon_be/mock"
	"coupon_be/util/logger"
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

	couponService, err := NewService(repoMock, writeDB)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, couponService)
	assert.Implements(t, (*Service)(nil), couponService)
}

type CouponServiceTestSuite struct {
	suite.Suite
	repo    *m.MockRepository
	writeDB *gorm.DB
	sqlMock sqlmock.Sqlmock
	ctx     context.Context

	couponService Service
}

func (suite *CouponServiceTestSuite) Before(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var err error

	logger.Initialise()

	suite.ctx = context.Background()
	suite.repo = m.NewMockRepository(ctrl)
	suite.writeDB, suite.sqlMock, err = m.NewMockDB()
	if err != nil {
		t.Fatal(err)
	}

	suite.couponService, err = NewService(suite.repo, suite.writeDB)
	if err != nil {
		t.Fatal(err)
	}
}

func (suite *CouponServiceTestSuite) After(t *testing.T) {}

func TestSuiteRunCouponService(t *testing.T) {
	suite.Run(t, new(CouponServiceTestSuite))
}
