package search

import (
	"xCut/entity"
	"xCut/search"
)

type SearchService struct {
	searchRepo search.SearchRepository
}

func (s SearchService) GetByName(keyword string) ([]entity.Shop, error) {
	return s.searchRepo.GetByName(keyword)
}

func (s SearchService) GetByLocation(lang float64, lat float64) ([]entity.Shop, error) {
	return s.searchRepo.GetByLocation(lang,lat)
}

func NewSearchService(searchRepo search.SearchRepository) search.SearchService {
	return &SearchService{searchRepo: searchRepo}
}

