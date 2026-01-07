package auth

import (
	"base_project/domain"
	"base_project/domain/enums"
	"base_project/request"
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util/logger"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (b *base) Register(ctx context.Context, input *request.Register) (*response.User, error) {
	logger.Info(ctx, "Register with request: %v", input)

	usernameExists, err := b.repository.FindUserByUsername(ctx, input.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if usernameExists != nil {
		return nil, sharedErrs.NewBusinessValidationErr(
			"Register Failed. Username %s already exists.", input.Username)
	}

	emailExists, err := b.repository.FindUserByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if emailExists != nil {
		return nil, sharedErrs.NewBusinessValidationErr(
			"Register Failed. Email %s already exists.", input.Email)
	}

	// encrypt password
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user, err := b.repository.CreateUser(ctx, &domain.User{
		BaseModel: domain.BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username: input.Username,
		Email:    input.Email,
		FullName: input.FullName,
		Password: string(password),
		Role:     enums.USERRole,
	})
	if err != nil {
		return nil, err
	}

	return &response.User{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role.String(),
	}, nil
}

//func (b *base) register(ctx context.Context) error {
//	now := time.Now()
//
//	for i := 0; i <= 100; i++ {
//		password, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
//		if err != nil {
//			return err
//		}
//
//		_, err = b.repository.CreateUser(ctx, &domain.User{
//			BaseModel: domain.BaseModel{
//				CreatedAt: now,
//				UpdatedAt: now,
//			},
//			Username: fmt.Sprintf("user%d", i),
//			Email:    fmt.Sprintf("user%d@mail.com", i),
//			FullName: fmt.Sprintf("User %d", i),
//			Password: string(password),
//			Role:     enums.USERRole,
//		})
//		if err != nil {
//			return err
//		}
//	}
//
//	password, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//
//	_, err = b.repository.CreateUser(ctx, &domain.User{
//		BaseModel: domain.BaseModel{
//			CreatedAt: now,
//			UpdatedAt: now,
//		},
//		Username: "admin",
//		Email:    "admin@mail.com",
//		FullName: "Administrator",
//		Password: string(password),
//		Role:     enums.ADMINRole,
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
