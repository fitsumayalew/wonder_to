package shop

import (
	"errors"
	"github.com/jinzhu/gorm"
	"xCut/entity"
	"xCut/shop"
)

type MockShopRepo struct {
	conn *gorm.DB
}

func (mShopRepo *MockShopRepo) GetShops() ([]entity.Shop, []error) {
	Shp := []entity.Shop{entity.MockShop}
	return Shp, nil
}

func (mShopRepo *MockShopRepo) GetShopByUserID(userID uint) (*entity.Shop, []error) {
	Shp := entity.MockShop
	if userID== 1 {
		return &Shp, nil
	}
	return nil, []error{errors.New("Not found")}
}

func NewMockShopRepo(db *gorm.DB) shop.ShopRepository {
	return &MockShopRepo{conn: db}
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
	Shp := shop
	return Shp,nil
}