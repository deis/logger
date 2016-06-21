package storage

import (
	"fmt"
	"log"

	r "gopkg.in/redis.v3"
)

type redisAdapter struct {
	bufferSize  int
	redisClient *r.Client
}

// NewRedisStorageAdapter returns a pointer to a new instance of a redis-based storage.Adapter.
func NewRedisStorageAdapter(bufferSize int) (*redisAdapter, error) {
	if bufferSize <= 0 {
		return nil, fmt.Errorf("Invalid buffer size: %d", bufferSize)
	}
	cfg, err := parseConfig(appName)
	if err != nil {
		log.Fatalf("config error: %s: ", err)
	}
	if err != nil {
		return nil, err
	}
	return &redisAdapter{
		bufferSize: bufferSize,
		redisClient: r.NewClient(&r.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword, // "" == no password
			DB:       int64(cfg.RedisDB),
		}),
	}, nil
}

// Write adds a log message to to an app-specific list in redis using ring-buffer-like semantics
func (a *redisAdapter) Write(app string, message string) error {
	// Note: Deliberately NOT using MULTI / transactions here since in this implementation of the
	// redis client, MULTI is not safe for concurrent use by multiple goroutines. It's been advised
	// by the authors of the gopkg.in/redis.v3 package to just use pipelining when possible...
	// and here that is technically possible. In the WORST case scenario, not having transactions
	// means we may momentarily have more than the desired number of log entries in the list /
	// buffer, but an LTRIM will eventually correct that, bringing the list / buffer back down to
	// its desired max size.
	pipeline := a.redisClient.Pipeline()
	if err := pipeline.RPush(app, message).Err(); err != nil {
		return err
	}
	if err := pipeline.LTrim(app, int64(-1*a.bufferSize), -1).Err(); err != nil {
		return err
	}
	if _, err := pipeline.Exec(); err != nil {
		return err
	}
	return nil
}

// Read retrieves a specified number of log lines from an app-specific list in redis
func (a *redisAdapter) Read(app string, lines int) ([]string, error) {
	stringSliceCmd := a.redisClient.LRange(app, int64(-1*lines), -1)
	result, err := stringSliceCmd.Result()
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	return nil, fmt.Errorf("Could not find logs for '%s'", app)
}

// Destroy deletes an app-specific list from redis
func (a *redisAdapter) Destroy(app string) error {
	if err := a.redisClient.Del(app).Err(); err != nil {
		return err
	}
	return nil
}

func (a *redisAdapter) Reopen() error {
	// No-op
	return nil
}
