package mock

import (
	"coupon_be/domain"
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
This file provides functionality to create instances of the specified required structs for unit testing purposes.
This ensures that tests have consistent and predictable data without the need for creating these objects manually in each test case.
*/

/*
 * ============================= MOCKING =============================
 */

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("Error occurs when opening a stub database connection : %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	}), &gorm.Config{})

	if err != nil {
		return nil, nil, fmt.Errorf("Error occurs when opening gorm database : %v", err)
	}

	return gormDB, mock, err
}

/*
 * ============================= DOMAIN =============================
 */

func InitUserDomain() *domain.User {
	now := time.Now()

	return &domain.User{
		BaseModel: domain.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username: "username",
		FullName: "full_name",
		Password: "password",
	}
}

func InitCouponDomain() *domain.Coupon {
	now := time.Now()

	return &domain.Coupon{
		BaseModel: domain.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:            "COUPON_TEST",
		Amount:          50,
		RemainingAmount: 50,
	}
}

func InitUserClaimDomain() *domain.UserClaim {
	now := time.Now()

	return &domain.UserClaim{
		BaseModel: domain.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:   1,
		CouponID: 1,
	}
}
