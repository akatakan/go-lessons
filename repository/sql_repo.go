package repository

import (
	"go-backend/models"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByID(id int) (models.User, bool)
	Save(user models.User)
}

type SQLUserRepository struct {
	db *gorm.DB
}

func NewSQLUserRepository(database *gorm.DB) *SQLUserRepository {
	return &SQLUserRepository{
		db: database,
	}
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
