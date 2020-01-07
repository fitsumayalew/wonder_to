package shop

import (
    "github.com/jinzhu/gorm"
	"xcut/entity"
	"xcut/shop"
)

type ShopGormRepo struct {
	conn *gorm.DB
}


func NewShopGormRepo(db *gorm.DB) shop.ShopRepository {
	return &ShopGormRepo{conn: db}
}



//  GetReviews returns all reviews stored in the database
func (shopRepo *ShopGormRepo) GetShops() ([]entity.Shop, []error) {
	shps := []entity.Shop{}
	errs := shopRepo.conn.Find(&shps).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return shps, errs
}

func (shopRepo *ShopGormRepo) GetShop(id uint) (*entity.Shop, []error) {
	shop := entity.Shop{}
	errs := shopRepo.conn.First(&shop,id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &shop, errs
}
// GetReviewByShopID retrieve a review from the database by its shopid

// UpdateReview updates a given review in the database
func (shopRepo *ShopGormRepo) UpdateShop(shop *entity.Shop) (*entity.Shop, []error){
	shps := shop
	errs := shopRepo.conn.Save(shps).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return shps, errs
}

// DeleteReview deletes a given review from the database
func (shopRepo *ShopGormRepo) DeleteShop(id uint) (*entity.Shop, []error) {
	shps, errs := shopRepo.GetShop(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = shopRepo.conn.Delete(shps, shps.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return shps, errs
}

// StoreReview stores a given review in the database
func (shopRepo *ShopGormRepo)StoreShop(shop *entity.Shop) (*entity.Shop, []error){
	shps := shop
	errs := shopRepo.conn.Create(shps).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return shps, errs
}

func (shopRepo *ShopGormRepo) GetShopByUserID(userID uint) (*entity.Shop, []error){
	shop := entity.Shop{}
	errs := shopRepo.conn.Set("gorm:auto_preload",true).Find(&shop, "user_id=?",userID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &shop, errs
}