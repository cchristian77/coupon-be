package redis

import (
	sharedErrs "coupon_be/shared/errors"
	"time"
)

const (
	retryDelay   = 500 * time.Millisecond
	expiry       = 5 * time.Second
	retriesCount = 3
)

type lockOptions struct {
	expiry       time.Duration
	retriesCount int
	retryDelay   time.Duration
}

type LockOption func(options *lockOptions) error

// SetLockExpiry expiry time for redis lock for the key
func SetLockExpiry(expiryTime time.Duration) LockOption {
	return func(options *lockOptions) error {
		if expiryTime < 1 {
			return sharedErrs.New(sharedErrs.ErrKindRedis, "expiry is less than 1")
		}

		options.expiry = expiryTime

		return nil
	}
}

// SetLockRetriesCount number of time, application tries to acquire lock.
func SetLockRetriesCount(retriesCount int) LockOption {
	return func(options *lockOptions) error {
		if retriesCount < 0 {
			return sharedErrs.New(sharedErrs.ErrKindRedis, "retries count is less than 1")
		}

		options.retriesCount = retriesCount

		return nil
	}
}

// SetLockRetryDelay time delay between 2 retries.
func SetLockRetryDelay(retryDelay time.Duration) LockOption {
	return func(options *lockOptions) error {
		if retryDelay <= 1 {
			return sharedErrs.New(sharedErrs.ErrKindRedis, "retryDelay is less than 1")
		}

		options.retryDelay = retryDelay

		return nil
	}
}

func getLockOptions(opts ...LockOption) (*lockOptions, error) {
	options := &lockOptions{
		expiry:       expiry,
		retriesCount: retriesCount,
		retryDelay:   retryDelay,
	}

	for _, o := range opts {
		if o != nil {
			if err := o(options); err != nil {
				return nil, sharedErrs.Wrap(err, "Redis lock option")
			}
		}
	}

	return options, nil
}

func overwriteLockOptions(globalOptions *lockOptions, opts ...LockOption) (*lockOptions, error) {
	if len(opts) == 0 {
		return globalOptions, nil
	}

	options := &lockOptions{
		expiry:       globalOptions.expiry,
		retriesCount: globalOptions.retriesCount,
		retryDelay:   globalOptions.retryDelay,
	}

	for _, o := range opts {
		if o != nil {
			if err := o(options); err != nil {
				return nil, sharedErrs.Wrap(err, "Redis lock option")
			}
		}
	}

	return options, nil
}
