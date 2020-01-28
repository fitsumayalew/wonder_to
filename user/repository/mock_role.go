package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"xCut/entity"
	"xCut/user"
)

type MockRoleRepo struct {
	conn *gorm.DB
}

func (mRoleRepo *MockRoleRepo) Roles() ([]entity.Role, []error) {
	return []entity.Role{entity.MockRole},nil
}

func (mRoleRepo *MockRoleRepo) Role(id uint) (*entity.Role, []error) {
	return  &entity.MockRole,nil
}

func (mRoleRepo *MockRoleRepo) RoleByName(name string) (*entity.Role, []error) {
	if name == "ADMIN" {
		return &entity.MockRole,nil
	}else{
		return nil, nil
	}
}

func NewMockRoleRepo(db *gorm.DB) user.RoleRepository {
	return &MockRoleRepo{conn: db}
}

func (roleRepo *RoleGormRepo) Roles() ([]entity.Role, []error) {
	Rol := []entity.Role{entity.MockRole}
	return Rol, nil
}

func (mRoleRepo *MockRoleRepo) GetRole(id uint) (*entity.Role, []error) {
	Rol := entity.MockRole
	if id != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return &Rol, nil
}

func (mRoleRepo *MockRoleRepo) UpdateRole(role *entity.Role) (*entity.Role, []error) {
	Rol := entity.MockRole
	return &Rol, nil
}

func (mRoleRepo *MockRoleRepo) DeleteRole(id uint) (*entity.Role, []error) {
	Rol := entity.MockRole
	if id != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return &Rol, nil
}

// StoreRole stores a given Role in the database
func (mRoleRepo *MockRoleRepo) StoreRole(role *entity.Role) (*entity.Role, []error) {
	Rol := role
	return Rol, nil
}
