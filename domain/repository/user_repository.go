package repository

import (
	"context"

	entity "github.com/josephakayesi/go-cerbos-abac/domain/core/entity"
	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmailOrUsername(c context.Context, uc vo.UserCredentials) (*entity.User, error)
	Create(c context.Context, user entity.User) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRespository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) FindByEmailOrUsername(c context.Context, uc vo.UserCredentials) (*entity.User, error) {
	user := entity.User{}

	tx := u.db.Where("email = ? OR username = ?", uc.Email, uc.Username).First(&user)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func (u *userRepository) Create(c context.Context, user entity.User) (*entity.User, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}

	tx.Create(&user)
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}

	tx.Commit()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}
