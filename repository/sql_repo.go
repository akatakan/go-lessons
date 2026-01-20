package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-backend/models"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var ctx = context.Background()

type IUserRepository interface {
	GetUserByID(id int) (models.User, bool)
	GetUserByUsername(username string) (models.User, error)
	Save(user models.User)
	Update(user models.User)
}

type SQLUserRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewSQLUserRepository(database *gorm.DB) *SQLUserRepository {
	return &SQLUserRepository{
		db: database,
		redisClient: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}

func (r *SQLUserRepository) GetUserWithCache(id int) (models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		json.Unmarshal([]byte(val), &user)
		return user, nil
	}

	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return user, err
	}

	userJson, _ := json.Marshal(user)
	r.redisClient.Set(ctx, cacheKey, userJson, 10*time.Minute)

	return user, nil
}

func (r *SQLUserRepository) GetUserByID(id int) (models.User, bool) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return models.User{}, false
	}
	return user, true
}

func (r *SQLUserRepository) Save(user models.User) {
	r.db.Create(&user)
}

func (r *SQLUserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := r.db.Where("kullanici_adi = ?", username).First(&user)
	if result.Error != nil {
		return models.User{}, errors.New("Kullanıcı bulunamadı")
	}
	return user, nil
}

func (r *SQLUserRepository) Update(user models.User) {
	r.db.Save(&user)
}
