package repository


import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/fitsumayalew/xCut/user"
	"github.com/fitsumayalew/xCut/entity"
)

type MockRoleRepo struct {
	conn *gorm.DB
}

func NewMockRoleRepo(db *gorm.DB) Role.RoleRepository {
	return &MockRoleRepo{conn: db}
}

func (roleRepo *RoleGormRepo) Roles() ([]entity.Role, []error) {
Rol := []entity.Role{entity.MockRole}
return Rol, nil
}


//func (roleRepo *RoleGormRepo) RoleByName(name string) (*entity.Role, []error) {
//	role := entity.Role{}
//	errs := roleRepo.conn.Find(&role, "name=?", name).GetErrors()
//	return &role, errs
//}

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
Rol := Role
return Rol,nil
}