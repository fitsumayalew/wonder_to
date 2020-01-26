package shop

import "xCut/entity"

type ShopRepository interface {
	GetShops() ([]entity.Shop, []error)
	GetShop(id uint) (*entity.Shop, []error)
	StoreShop(shop *entity.Shop) (*entity.Shop, []error)
	UpdateShop(shop *entity.Shop) (*entity.Shop, []error)
	DeleteShop(id uint) (*entity.Shop, []error)
	GetShopByUserID(userID uint) (*entity.Shop, []error)
}

