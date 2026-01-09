package redis

import (
	"context"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/util/logger"

	"github.com/go-redsync/redsync"
)

const (
	redisLockExpiry       = 20  // in seconds
	redisLockRetryDelay   = 300 // in miliseconds
	redisLockMaxNoOfRetry = 3
)

// Lock RedisLock - redis implementation for distributed locking
type Lock struct {
	rSync   *redsync.Redsync
	options *lockOptions
}

// ILock - interface for distributed locking
type ILock interface {
	WithLock(ctx context.Context, key string, fn func() error, options ...LockOption) error
}

func getLock(key string, rSync *redsync.Redsync, option *lockOptions) (*redsync.Mutex, error) {
	if rSync == nil {
		return nil, sharedErrs.New(sharedErrs.ErrKindRedis, "Redis pool is nil")
	}

	var lockOptions []redsync.Option

	expiryOpt := redsync.SetExpiry(option.expiry)
	triesOpt := redsync.SetTries(option.retriesCount)
	retryDelayOpt := redsync.SetRetryDelay(option.retryDelay)
	lockOptions = append(lockOptions, expiryOpt, triesOpt, retryDelayOpt)

	mutex := rSync.NewMutex(key, lockOptions...)

	if err := mutex.Lock(); err != nil {
		return nil, sharedErrs.New(sharedErrs.ErrKindAcquireRedisLock, "Error acquiring redis lock: %v", err)
	}

	return mutex, nil
}

// WithLock - Wrapper locking function for redis
func (l *Lock) WithLock(ctx context.Context, key string, fn func() error, options ...LockOption) error {
	if l == nil {
		return sharedErrs.New(sharedErrs.ErrKindRedis, "Lock object is nil")
	}

	opts, err := overwriteLockOptions(l.options, options...)
	if err != nil {
		return sharedErrs.Wrap(err, "WithLock applying option")
	}

	mutex, err := getLock(key, l.rSync, opts)
	if err != nil {
		return sharedErrs.New(sharedErrs.ErrKindAcquireRedisLock, "Error acquiring redis lock: %v", err)
	}

	defer func() {
		_, err = mutex.Unlock()
		if err != nil {
			logger.Error(ctx, "Error while unlock %v", err)
		}

		logger.Debug(ctx, "Lock released successfully for key %v", key)
	}()

	logger.Debug(ctx, "Acquired lock for key %v", key)

	if err = fn(); err != nil {
		return err
	}

	return nil
}

// newLock Provider/Factory function to return a redis lock struct
func newLock(ctx context.Context, pool *Pool, options ...LockOption) (ILock, error) {
	if pool == nil {
		return nil, sharedErrs.New(sharedErrs.ErrKindRedis, "Redis pool is nil")
	}

	opts, err := getLockOptions(options...)
	if err != nil {
		return nil, sharedErrs.Wrap(err, "initialising lock")
	}

	logger.Debug(ctx, "Redis lock initialised successfully")

	return &Lock{
		rSync:   redsync.New([]redsync.Pool{pool.Pool}),
		options: opts,
	}, nil
}
