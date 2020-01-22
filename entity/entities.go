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
	Price         float32
	EstimatedTime uint
}

type User struct {
	gorm.Model
	FullName string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null;unique"`
	Password string `gorm:"type:varchar(255);not null"`
	Phone    string `gorm:"type:varchar(16);not null; unique"`
	RoleID   uint
}

//Shop Represents the the shops the user views
type Shop struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255);not null"`
	City      string
	Lat       float64
	Long      float64
	Address   string `gorm:"type:varchar(255);not null; unique"`
	Phone     string `gorm:"type:varchar(32)"`
	Website   string `gorm:"type:varchar(255)"`
	Image     string `gorm:"type:varchar(255)"`
	UserID    uint
	WeekDayOpenHour uint `gorm:"default:480"`
	WeekDayCloseHour uint `gorm:"default:1200"`
	WeekendOpenHour uint `gorm:"default:510"`
	WeekendCloseHour uint `gorm:"default:1080"`
	Services  []Service
}


type Appointment struct {
	gorm.Model
	UserID     uint
	ShopID     uint
	ServicesID uint
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
	Review string `gorm:"type:varchar(1024);not null"`
	Reply  string `gorm:"type:varchar(1024);not null"`
	Rating float32 `gorm:"not null"`
}