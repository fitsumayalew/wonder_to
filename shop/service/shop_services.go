package shop

import (
	"xcut/entity"
	"xcut/shop"
)

type ShopService struct {
	shopRepo shop.ShopRepository
}

func (sp *ShopService) GetShopByUserID(userID uint) (*entity.Shop, []error) {
	shp, errs := sp.shopRepo.GetShopByUserID(userID)
	if len(errs) > 0 {
		return nil, errs
	}
	return shp, errs
}

func NewShopService(shpRepo shop.ShopRepository) shop.ShopService {
	return &ShopService{shopRepo: shpRepo}
}

// GetReviews returns all stored reviews
func (sp *ShopService) GetShops() ([]entity.Shop, []error) {
	shp, errs := sp.shopRepo.GetShops()
	if len(errs) > 0 {
		return nil, errs
	}
	return shp, errs
}

func (sp *ShopService) GetShop(id uint) (*entity.Shop, []error) {
	shop, errs := sp.shopRepo.GetShop(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return shop, errs
}

// UpdateReview updates a given review in a database
func (sp *ShopService) UpdateShop(shop *entity.Shop) (*entity.Shop, []error) {
	shp, errs := sp.shopRepo.UpdateShop(shop)
	if len(errs) > 0 {
		return nil, errs
	}
	return shp, errs
}

// DeleteReview deletes a given review
func (sp *ShopService) DeleteShop(id uint) (*entity.Shop, []error) {
	shp, errs := sp.shopRepo.DeleteShop(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return shp, errs
}

// StoreReview stores a given review
func (sp *ShopService) StoreShop(shop *entity.Shop) (*entity.Shop, []error) {
	shp, errs := sp.shopRepo.StoreShop(shop)
	if len(errs) > 0 {
		return nil, errs
	}
	return shp, errs
}
