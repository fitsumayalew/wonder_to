package search

import (
	"github.com/jinzhu/gorm"
	"xCut/entity"
	"xCut/search"
)

type SearchGormRepo struct {
	conn *gorm.DB
}

func (s SearchGormRepo) GetByName(keyword string) ([]entity.Shop, error) {
	shops := []entity.Shop{}
	s.conn.Where("name like ?", "%"+keyword+"%").Find(&shops)
	return shops, nil
}

func (s SearchGormRepo) GetByLocation(lang float64, lat float64) ([]entity.Shop, error) {
	shops := []entity.Shop{}
	s.conn.Where("lat > ? and lat < ? and long > ? and long < ?", lat-0.05,lat+0.05,lang-0.05,lang+0.05).Find(&shops)
	return shops, nil
}

func NewSearchGormRepo(db *gorm.DB) search.SearchRepository {
	return &SearchGormRepo{conn: db}
}
