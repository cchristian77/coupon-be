package user

import (
	"context"
	"coupon_be/domain"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (b *base) Register(ctx context.Context) error {
	now := time.Now()

	for i := 0; i <= 100; i++ {
		password, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = b.repository.CreateUser(ctx, &domain.User{
			BaseModel: domain.BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username: fmt.Sprintf("user_%d", i),
			FullName: fmt.Sprintf("User %d", i),
			Password: string(password),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
