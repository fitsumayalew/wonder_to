package entity

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Role struct {
	ID   uint
	Name string `gorm:"type:varchar(255)"`
}

type Service struct {
	gorm.Model
	Name          string `gorm:"type:varchar(255);not null"`
	ShopID        uint
	Shop          Shop
	Price         float32
	EstimatedTime uint
	Image         string `gorm:"type:varchar(128);not null`
}

type User struct {
	gorm.Model
	FullName        string `gorm:"type:varchar(255);not null"`
	Email           string `gorm:"type:varchar(255);not null;unique"`
	Password        string `gorm:"type:varchar(255);not null"`
	Phone           string `gorm:"type:varchar(16);not null; unique"`
	Image           string `gorm:"type:varchar(128);`
	NumberOfReviews uint
	RoleID          uint
}



type Shop struct {
	gorm.Model
	Name             string  `json:"name" gorm:"type:varchar(255);not null"`
	City             string  `json:"city"`
	Lat              float64 `json:"lat"`
	Long             float64 `json:"long"`
	Address          string  `json:"address" gorm:"type:varchar(255);not null;"`
	Phone            string  `json:"phone" gorm:"type:varchar(32)"`
	Website          string  `json:"website" gorm:"type:varchar(255)"`
	Image            string  `json:"image"  gorm:"type:varchar(255)"`
	UserID           uint    `json:"userid"`
	WeekDayOpenHour  uint    `json:"dayopen" gorm:"default:480"`
	WeekDayCloseHour uint    `json:"dayclose" gorm:"default:1200"`
	WeekendOpenHour  uint    `json:"open" gorm:"default:510"`
	WeekendCloseHour uint    `json:"close" gorm:"default:1080"`
	Rating           uint    `json:"rating" gorm:"default:0"`
}

type Appointment struct {
	gorm.Model
	UserID          uint
	ShopID          uint
	ServicesID      uint
	User            User
	Shop            Shop
	AppointmentTime *time.Time
}

type Session struct {
	gorm.Model
	SessionId  string `gorm:"type:varchar(255);not null"`
	UUID       uint
	Expires    int64  `gorm:"type:varchar(255);not null"`
	SigningKey []byte `gorm:"type:varchar(255);not null"`
}

type Review struct {
	gorm.Model
	UserID uint
	ShopID uint
	User   User
	Shop   Shop
	Review string `gorm:"type:varchar(1024);not null"`
	Reply  string `gorm:"type:varchar(1024);not null"`
	Rating uint   `gorm:"not null"`
}
