package search

import (
	"github.com/jinzhu/gorm"
	"xCut/entity"
	"xCut/review"
)

type SearchGormRepo struct {
	conn *gorm.DB
}

func NewSearchGormRepo(db *gorm.DB) search. {
	return &SearchGormRepo{conn: db}
}
