// user/repository/redis_mock.go
package repository

import (
	"fmt"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"log"
)

type redisMockRepository struct {
	// In-memory storage for testing
	cache map[string]string
}

func NewRedisMockRepository() domain.UserCacheRepository {
	return &redisMockRepository{
		cache: make(map[string]string),
	}
}

func (r *redisMockRepository) SetUserCache(userID uint, user *entity.User) error {
	log.Printf("Mock Redis: Setting user cache for ID %d", userID)
	// In real implementation, this would serialize user to JSON
	r.cache[fmt.Sprintf("user:%d", userID)] = "cached_user_data"
	return nil
}

func (r *redisMockRepository) GetUserCache(userID uint) (*entity.User, error) {
	log.Printf("Mock Redis: Getting user cache for ID %d", userID)
	key := fmt.Sprintf("user:%d", userID)
	if _, exists := r.cache[key]; exists {
		log.Printf("Mock Redis: Cache hit for user ID %d", userID)
		// In real implementation, this would deserialize from JSON
		return nil, nil // Return nil to indicate cache miss for now
	}
	log.Printf("Mock Redis: Cache miss for user ID %d", userID)
	return nil, nil
}

func (r *redisMockRepository) DeleteUserCache(userID uint) error {
	log.Printf("Mock Redis: Deleting user cache for ID %d", userID)
	key := fmt.Sprintf("user:%d", userID)
	delete(r.cache, key)
	return nil
}

func (r *redisMockRepository) SetUserSession(sessionID string, userID uint) error {
	log.Printf("Mock Redis: Setting session %s for user ID %d", sessionID, userID)
	r.cache[fmt.Sprintf("session:%s", sessionID)] = fmt.Sprintf("%d", userID)
	return nil
}

func (r *redisMockRepository) GetUserSession(sessionID string) (uint, error) {
	log.Printf("Mock Redis: Getting session %s", sessionID)
	key := fmt.Sprintf("session:%s", sessionID)
	if value, exists := r.cache[key]; exists {
		log.Printf("Mock Redis: Session found: %s", value)
		// Parse userID from string
		var userID uint
		fmt.Sscanf(value, "%d", &userID)
		return userID, nil
	}
	log.Printf("Mock Redis: Session not found")
	return 0, nil
}
