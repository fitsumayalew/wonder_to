package service

import (
	"xCut/entity"
	"xCut/review"
)

type ReviewService struct {
	reviewRepo review.ReviewRepository
}



func NewReviewService(revoRepo review.ReviewRepository) review.ReviewService {
	return &ReviewService{reviewRepo: revoRepo}
}

// GetReviews returns all stored reviews
func (rs *ReviewService) GetReviews() ([]entity.Review, []error) {
	revs, errs := rs.reviewRepo.GetReviews()
	if len(errs) > 0 {
		return nil, errs
	}
	return revs, errs
}

func (rs *ReviewService) GetReview(id uint) (*entity.Review, []error) {
	rev, errs := rs.reviewRepo.GetReview(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return rev, errs
}

// GetReviewByShopID  retrieves stored comment by its id
func (rs *ReviewService) GetReviewsByShopID(ShopID uint) ([]entity.Review, []error) {
	revs, errs := rs.reviewRepo.GetReviewsByShopID(ShopID)
	if len(errs) > 0 {
		return nil, errs
	}
	return revs, errs
}

// UpdateReview updates a given review in a database
func (rs *ReviewService) UpdateReview(review *entity.Review) (*entity.Review, []error) {
	revs, errs := rs.reviewRepo.UpdateReview(review)
	if len(errs) > 0 {
		return nil, errs
	}
	return revs, errs
}

// DeleteReview deletes a given review
func (rs *ReviewService) DeleteReview(id uint) (*entity.Review, []error) {
	revs, errs := rs.reviewRepo.DeleteReview(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return revs, errs
}

// StoreReview stores a given review
func (rs *ReviewService) StoreReview(review *entity.Review) (*entity.Review, []error) {
	revs, errs := rs.reviewRepo.StoreReview(review)
	if len(errs) > 0 {
		return nil, errs
	}
	return revs, errs
}
