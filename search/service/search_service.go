package search

import (
	"xCut/review"
)

type LocationService struct {
	locationRepo review.ReviewRepository
}

func LocationService(revoRepo review.ReviewRepository) review.ReviewService {
	return &ReviewService{reviewRepo: revoRepo}
}

