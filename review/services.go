package review

import (
	"xcut/entity"
)

type ReviewService interface{
	
	GetReviews() ([]entity.Review, []error)
	GetReview(id uint) (*entity.Review, []error)
		//getRecentReview
	StoreReview(review *entity.Review) (*entity.Review, []error)
	UpdateReview(review *entity.Review) (*entity.Review, []error)
	DeleteReview(id uint) (*entity.Review, []error)
	GetReviewsByShopID(ShopID uint) ([]entity.Review, []error)

}