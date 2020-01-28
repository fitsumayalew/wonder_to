package repository

import (
	"github.com/jinzhu/gorm"
	"xCut/entity"
	"xCut/user"
)

type MockUserRepo struct {
	conn *gorm.DB
}

func (userRepo *MockUserRepo) User(id uint) (*entity.User, []error) {
	users := entity.MockUSer
	return &users, nil
}

func (userRepo *MockUserRepo) UserByEmail(email string) (*entity.User, []error) {
	users := entity.MockUSer
	return &users, nil
}

func NewMockUserRepo(db *gorm.DB) user.UserRepository {
	return &MockUserRepo{conn: db}
}

func (userRepo *MockUserRepo) GetUser() ([]entity.User, []error) {
	users := []entity.User{entity.MockUSer}
	return users, nil
}

func (userRepo *MockUserRepo) StoreUser(user *entity.User) (*entity.User, []error) {
	users := user
	return users, nil
}
