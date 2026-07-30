package repository

import (
	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

// User is the GORM model backing the `users` table.
// TODO: 실제 스키마에 맞게 컬럼을 조정하세요.
type User struct {
	ID    uint64 `gorm:"primaryKey;column:id"`
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email;uniqueIndex"`
}

func (User) TableName() string {
	return "users"
}

// UserRepository is the persistence boundary for the User aggregate.
// Controllers/services must depend on this interface, never on *gorm.DB directly.
type UserRepository interface {
	FindByID(id uint64) (*User, error)
	Create(user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uint64) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find user by id")
	}
	return &user, nil
}

func (r *userRepository) Create(user *User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}
