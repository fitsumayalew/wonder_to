package search

import (
	"xCut/entity"
)

type SearchService interface {
	GetByName(keyword string) ([]entity.Shop, error)
	GetByLocation(lang float64, lat float64) ([]entity.Shop, error)
}