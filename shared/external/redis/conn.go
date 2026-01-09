package redis

import (
	"context"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/util/logger"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	connectionType               = "tcp"
	defaultConnectTimeoutSeconds = 5
	defaultReadTimeoutSeconds    = 1
	defaultWriteTimeoutSeconds   = 1
)

var (
	defaultUseTLS        = true
	defaultTLSSkipVerify = true
)

type Pool struct {
	*redis.Pool
}

// ConnectRedis creates and returns a new Redis connection with the given configuration.
func connectRedis(ctx context.Context, config *config) (*Pool, error) {
	if config.ConnectTimeoutSeconds == 0 {
		config.ConnectTimeoutSeconds = defaultConnectTimeoutSeconds
	}
	if config.ReadTimeoutSeconds == 0 {
		config.ReadTimeoutSeconds = defaultReadTimeoutSeconds
	}
	if config.WriteTimeoutSeconds == 0 {
		config.WriteTimeoutSeconds = defaultWriteTimeoutSeconds
	}
	if config.UseTLS == nil {
		config.UseTLS = &defaultUseTLS
	}
	if config.TLSSkipVerify == nil {
		config.TLSSkipVerify = &defaultTLSSkipVerify
	}

	pool := &redis.Pool{
		MaxIdle:     config.MaxIdleConnections,
		MaxActive:   config.MaxActiveConnections,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(connectionType, fmt.Sprintf("%v:%v", config.Host, config.Port),
				redis.DialPassword(config.Password),
				redis.DialUseTLS(*config.UseTLS),
				redis.DialTLSSkipVerify(*config.TLSSkipVerify),
				redis.DialConnectTimeout(time.Duration(config.ConnectTimeoutSeconds)*time.Second),
				redis.DialReadTimeout(time.Duration(config.ReadTimeoutSeconds)*time.Second),
				redis.DialWriteTimeout(time.Duration(config.WriteTimeoutSeconds)*time.Second),
			)
			if err != nil {
				msg := fmt.Sprintf("Error connecting to redis at %v:%v", config.Host, config.Port)
				return nil, sharedErrs.NewWithCause(sharedErrs.ErrKindRedis, msg, err)
			}
			return c, err
		},
	}

	conn := &Pool{pool}
	if _, err := conn.ping(); err != nil {
		return nil, sharedErrs.Wrap(err, "redis ping")
	}

	logger.Info(ctx, "Redis connected successfully at %v:%v", config.Host, config.Port)

	return conn, nil
}

// Ping is a utility to ping a Redis server to verify that the connection is created successfully.
func (r *Pool) ping() (string, error) {
	if r == nil {
		return "", sharedErrs.New(sharedErrs.ErrKindRedis, "Redis pool is nil")
	}

	conn, err := r.Dial()
	if err != nil {
		return "", sharedErrs.NewWithCause(sharedErrs.ErrKindRedis, "failed to dial to redis", err)
	}

	resp, err := conn.Do("ping")
	if err != nil {
		return "", sharedErrs.NewWithCause(sharedErrs.ErrKindRedis, "failed to ping redis", err)
	}

	return resp.(string), nil
}
