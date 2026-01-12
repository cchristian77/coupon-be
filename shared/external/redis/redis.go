package redis

import (
	"context"
	sharedErrs "coupon_be/shared/errors"
	uConfig "coupon_be/util/config"
	"time"
)

var redisClient *Pool

// GetConnection - Returns the redis connection
func GetConnection(ctx context.Context) (*Pool, error) {
	if redisClient == nil {
		redisConfig := uConfig.Env().Redis

		pool, err := connectRedis(ctx, &config{
			Host:                 redisConfig.Host,
			Port:                 redisConfig.Port,
			Password:             redisConfig.Password,
			MaxIdleConnections:   redisConfig.MaxIdleConnections,
			MaxActiveConnections: redisConfig.MaxActiveConnections,
			IdleTimeout:          redisConfig.IdleTimeout,
			UseTLS:               &redisConfig.UseTLS,
		})
		if err != nil {
			return nil, err
		}

		redisClient = pool
	}

	return redisClient, nil
}

func GetRedisLock(ctx context.Context) (ILock, error) {
	redis, err := GetConnection(ctx)
	if err != nil {
		return nil, sharedErrs.Wrap(err, "get redis conn")
	}

	return newLock(ctx, redis,
		SetLockRetriesCount(redisLockMaxNoOfRetry),
		SetLockRetryDelay(redisLockRetryDelay*time.Millisecond),
		SetLockExpiry(redisLockExpiry*time.Second),
	)
}
