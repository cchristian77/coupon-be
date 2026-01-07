package mock

import (
	"base_project/domain"
	"base_project/domain/enums"
	"base_project/response"
	"base_project/util"
	tokenMaker "base_project/util/token"
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
This file provides functionality to create instances of the specified required structs for unit testing purposes.
This ensures that tests have consistent and predictable data without the need for creating these objects manually in each test case.
*/

const DefaultJWTSecretForTest = "secret"

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

func MockAccessToken() (string, *domain.Session) {
	tokenMaker.Initialise(DefaultJWTSecretForTest)

	sessionID := uuid.New()
	userID := uint64(1)
	now := time.Now()
	expiresAt := now.Add(1 * time.Hour)

	accessToken, _, _ := tokenMaker.Get().Generate(sessionID, userID, 1*time.Hour)

	session := &domain.Session{
		ID:                   1,
		SessionID:            sessionID,
		UserID:               userID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: expiresAt,
		AccessTokenCreatedAt: now,
	}

	return accessToken, session
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
		DeletedAt: nil,
		Username:  "username",
		Email:     "test@mail.com",
		FullName:  "full_name",
		Password:  "password",
		Role:      enums.ADMINRole,
	}
}

func InitCommentDomain() *domain.Comment {
	now := time.Now()

	return &domain.Comment{
		BaseModel: domain.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:  1,
		PostID:  1,
		Comment: "comment test",
		Rating:  util.ToPointerValue(uint8(5)),
	}
}

func InitPostDomain() *domain.Post {
	now := time.Now()

	return &domain.Post{
		BaseModel: domain.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID: 1,
		Slug:   "slug-test",
		Title:  "title test",
		Body:   "body test",
		Status: enums.PUBLISHEDPostStatus,
	}
}

/*
 * ============================= Response =============================
 */

func InitCommentResponse() *response.Comment {
	now := time.Now()

	return &response.Comment{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "username",
		Comment:   "comment test",
	}
}
