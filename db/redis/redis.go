package redis

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

var once sync.Once
var instance *Redis

func initialize(envType *string) *redis.Client {
	slog.Debug("Initializing Redis")
	var rdb *redis.Client
	if *envType == "dev" {
		rdb = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		})
	}
	return rdb
}

// ping returns
func (r *Redis) Ping() string {
	return fmt.Sprintf("PONG1 %s", r.client.Ping(context.Background()).Val())
}

func NewRedis(envType *string) *Redis {
	once.Do(func() {
		client := initialize(envType)
		instance = &Redis{
			client: client,
		}
		slog.Debug("Connected with Redis!!!!!")
		//call ping
		slog.Info(instance.Ping())
	})
	return instance
}

func (r *Redis) StoreEmailHash(email string) (string, error) {
	// create a sha256 hash of this email
	h := sha256.New()
	h.Write([]byte(email))
	bs := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// store hash:email
	return string(bs), r.client.Set(context.Background(), string(bs), email, 0).Err()
}

// fetch email from hash and delete the key
func (r *Redis) GetEmailFromHash(hash string) (string, error) {
	ctx := context.Background()
	email, err := r.client.Get(ctx, hash).Result()
	if err != nil || len(email) == 0 {
		return "", err
	}
	// now delete the key
	err = r.client.Del(ctx, hash).Err()
	return email, err
}

// GenerateToken generates a new token for password reset
func (r *Redis) GenerateToken(email string) (string, error) {
	// Create a SHA256 hash of the email for a secure token
	h := sha256.New()
	h.Write([]byte(email + time.Now().String())) // Include timestamp for uniqueness
	token := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Store the token in Redis with an expiration of 1 hour
	if err := r.client.Set(context.Background(), token, email, time.Hour).Err(); err != nil {
		return "", fmt.Errorf("failed to store token in Redis: %w", err)
	}

	return token, nil
}

// delete token

func (r *Redis) DeleteToken(token string) error {
	return r.client.Del(context.Background(), token).Err()
}
