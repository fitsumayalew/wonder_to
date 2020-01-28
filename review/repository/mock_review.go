package repository

import (
	"github.com/jinzhu/gorm"
	"xCut/entity"
	"xCut/review"
)

type MockReviewRepo struct {
	conn *gorm.DB
}

// NewReviewtGormRepo returns new object of MockReviewRepo
func NewMockReviewRepo( db *gorm.DB) review.ReviewRepository {
	return &MockReviewRepo{conn: db}
}

//  GetReviews returns all reviews stored in the database
func (mrevRepo *MockReviewRepo) GetReviews() ([]entity.Review, []error) {
	revo := []entity.Review{entity.MockReview}
	return revo, nil
}

// GetReviewByShopID retrieve a review from the database by its shopid
func (mrevRepo *MockReviewRepo) GetReviewsByShopID(ShopID uint) ([]entity.Review, []error) {
	revo := entity.MockReview{}
	if id == 1 {
		return &revo, nil
	}
	return nil, []error{errors.New("Not found")}
}
}

// UpdateReview updates a given review in the database
func (mrevRepo *MockReviewRepo) UpdateReview(review *entity.Review) (*entity.Review, []error) {
	revo := entity.MockReview
	return &revo, nil
}



// DeleteReview deletes a given review from the database
func (mrevRepo *MockReviewRepo) DeleteReview(id uint) (*entity.Review, []error) {
	revo := entity.MockReview
	if id != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return &revo, nil
}

// StoreReview stores a given review in the database
func (mrevRepo *MockReviewRepo) StoreReview(review *entity.Review) (*entity.Review, []error) {
	revo := review
return revo,nil
}
