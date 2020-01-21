package repository

import (
	"github.com/jinzhu/gorm"
	"xCut/entity"
	"xCut/review"
)

type ReviewGormRepo struct {
	conn *gorm.DB
}

// NewReviewtGormRepo returns new object of ReviewGormRepo
func NewReviewGormRepo(db *gorm.DB) review.ReviewRepository {
	return &ReviewGormRepo{conn: db}
}

//  GetReviews returns all reviews stored in the database
func (revRepo *ReviewGormRepo) GetReviews() ([]entity.Review, []error) {
	revws := []entity.Review{}
	errs := revRepo.conn.Find(&revws).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return revws, errs
}

// GetReviewByShopID retrieve a review from the database by its shopid
func (revRepo *ReviewGormRepo) GetReviewsByShopID(ShopID uint) ([]entity.Review, []error) {
	revo := []entity.Review{}
	errs := revRepo.conn.Set("gorm:auto_preload", true).Find(&revo, "shop_id=?", ShopID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return revo, errs
}

// UpdateReview updates a given review in the database
func (revRepo *ReviewGormRepo) UpdateReview(review *entity.Review) (*entity.Review, []error) {
	revi := review
	errs := revRepo.conn.Save(revi).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return revi, errs
}

func (revRepo *ReviewGormRepo) GetReview(id uint) (*entity.Review, []error) {
	review := entity.Review{}
	errs := revRepo.conn.Set("gorm:auto_preload", true).First(&review, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &review, errs
}

// DeleteReview deletes a given review from the database
func (revRepo *ReviewGormRepo) DeleteReview(id uint) (*entity.Review, []error) {
	review, errs := revRepo.GetReview(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = revRepo.conn.Delete(review, review.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return review, errs
}

// StoreReview stores a given review in the database
func (revRepo *ReviewGormRepo) StoreReview(review *entity.Review) (*entity.Review, []error) {
	revi := review
	errs := revRepo.conn.Create(revi).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return revi, errs
}
