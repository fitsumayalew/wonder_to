package repository

import (
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/fitsumayalew/xCut/shop"
	"github.com/fitsumayalew/xCut/entity"
)

type MockShopRepo struct {
	conn *gorm.DB
}



func NewMockShopRepo(db *gorm.DB) Shop.ShopRepository {
	return &MockShopRepo{conn: db}
}





// GetShopByShopID retrieve a Shop from the database by its shopid
func (mShopRepo *MockShopRepo) GetShopsByUserID(userID uint) ([]entity.Shop, []error) {
	Shp := entity.MockShop{}
	if userID== 1 {
		return &Shp, nil
	}
	return nil, []error{errors.New("Not found")}
}

// UpdateShop updates a given Shop in the database
func (mShopRepo *MockShopRepo) UpdateShop(shop *entity.Shop) (*entity.Shop, []error) {
	Shp := entity.MockShop
	return &Shp, nil
}

func (mShopRepo *MockShopRepo) GetShop(id uint) (*entity.Shop, []error) {
	Shp := entity.MockShop
	if id != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return &Shp, nil
}
// DeleteShop deletes a given Shop from the database
func (mShopRepo *MockShopRepo) DeleteShop(id uint) (*entity.Shop, []error) {
	Shp := entity.MockShop
	if id != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return &Shp, nil
}

// StoreShop stores a given Shop in the database
func (mShopRepo *MockShopRepo) StoreShop(shop *entity.Shop) (*entity.Shop, []error) {
	Shp := Shop
	return Shp,nil
}