package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/fitsumayalew/xCut/entity"
	"github.com/fitsumayalew/xCut/user"
)






type MockUserRepo struct {
	conn *gorm.DB
}


func NewMockUserRepo( db *gorm.DB) user.UserRepository {
	return &MockUserRepo{conn: db}
}


func (userRepo *MockUserRepo) GetUser() ([]entity.Review, []error) {
	users := []entity.Review{entity.MockUser}
	return users, nil
}


//func (userRepo *UserGormRepo) UserByEmail(email string) (*entity.User, []error) {
//	user := entity.User{}
//	errs := userRepo.conn.First(&user, "email=?", email).GetErrors()
//	return &user, errs
//}


func (userRepo *MockUserRepo) StoreUser(user *entity.User) (*entity.User, []error) {
	users := user
	return users,nil
}
