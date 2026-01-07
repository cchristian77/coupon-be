package repository

import (
	"base_project/domain"
	sharedErrs "base_project/shared/errors"
	"base_project/shared/external/database"
	"base_project/util/logger"
	"context"

	"gorm.io/gorm/clause"
)

func (r *repo) CreateSession(ctx context.Context, data *domain.Session) (*domain.Session, error) {
	err := r.DB.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on create session : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) FindSessionBySessionID(ctx context.Context, sessionID string) (*domain.Session, error) {
	var result *domain.Session

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find session by id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) DeleteSessionByID(ctx context.Context, id uint64) error {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Session{}).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on delete session by id : %v", err)

		return sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return nil
}
