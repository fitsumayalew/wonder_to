package entity

import (
	"github.com/jinzhu/gorm"
)

type Role struct {
	ID   uint
	Name string `gorm:"type:varchar(255)"`
}

type Services struct {
	gorm.Model
	Name          string `gorm:"type:varchar(255);not null"`
	ShopID        uint
	Price         float32
	EstimatedTime uint
}

type User struct {
	gorm.Model
	FullName string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null;unique"`
	Password string `gorm:"type:varchar(255)"`
	Phone    string `gorm:"type:varchar(16);not null; unique"`
	RoleID   uint
}

//Shop Represents the the shops the user views
type Shop struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null"`
	City     string
	Lat      float32
	Long     float32
	Address string `gorm:"type:varchar(255);not null; unique"`
	Phone    string  `gorm:"type:varchar(32);not null; unique"`
	Website  string `gorm:"type:varchar(255);not null; unique"`
	Image    string `gorm:"type:varchar(255)"`
	UserID	uint
	Services []Services
}


type Appointments struct {
	gorm.Model
	UserID     uint
	ShopID     uint
	ServicesID uint
}

type Session struct {
	gorm.Model
	SessionId  string `gorm:"type:varchar(255);not null"`
	UUID       uint
	Expires    int64  `gorm:"type:varchar(255);not null"`
	SigningKey []byte `gorm:"type:varchar(255);not null"`
}






