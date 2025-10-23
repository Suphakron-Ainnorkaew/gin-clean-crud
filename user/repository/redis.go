// user/repository/redis.go
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisUserRepository struct {
	client *redis.Client
}

func NewRedisUserRepository(client *redis.Client) domain.UserCacheRepository {
	return &redisUserRepository{client: client}
}

func (r *redisUserRepository) SetUserCache(userID uint, user *entity.User) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:%d", userID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, 1*time.Hour).Err()
}

func (r *redisUserRepository) GetUserCache(userID uint) (*entity.User, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%d", userID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var user entity.User
	err = json.Unmarshal([]byte(data), &user)
	return &user, err
}

func (r *redisUserRepository) DeleteUserCache(userID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:%d", userID)
	return r.client.Del(ctx, key).Err()
}

func (r *redisUserRepository) SetUserSession(sessionID string, userID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", sessionID)
	return r.client.Set(ctx, key, userID, 24*time.Hour).Err()
}

func (r *redisUserRepository) GetUserSession(sessionID string) (uint, error) {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", sessionID)
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	
	var userID uint
	_, err = fmt.Sscanf(result, "%d", &userID)
	return userID, err
}