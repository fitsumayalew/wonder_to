package entity

import (
	"github.com/jinzhu/gorm"
	"time"
)

var MockRole = Role{
	ID:  1,
	Name:  "Mock Role 01",
}

var MockService = Service{
  	Model:gorm.Model{
  		ID:1,
		CreatedAt: time.Time{},
  	},
	Name:"Mock Service",
	ShopID:1,
	Shop:MockShop,
	Price:75,
	Image:"mock_services.png",
	EstimatedTime: 10,

}

var MockUSer = User{
	Model:gorm.Model{
		ID:1,
		CreatedAt: time.Time{},
	},
	FullName:"Mock User 01",
	Email:"mockuser@gmail.com",
	Password:"P@$$w0rd",
	Phone:"0901010101",
	NumberOfReviews:3,
	Image:"mock_user.png",
	RoleID: 1,
	}

//Shop Represents the the shops the user views
var MockShop = Shop{
	Model:gorm.Model{
		ID:1,
		CreatedAt: time.Time{},
	},
	Name:"Mock Shop 01",
	City:"Mock City 01",
	Lat:10.2,
	Long:10.5,
	Address:"Mock Adress 01",
	Phone:"0901010101",
	Website:"WWW.mocksite.com",
	Image:"mock_item.png",
	UserID:1,
	WeekDayOpenHour:2,
	WeekDayCloseHour:3,
	WeekendOpenHour:2,
	WeekendCloseHour:3,
}


var MockAppointment =Appointment{
	Model:gorm.Model{
		ID:1,
		CreatedAt: time.Time{},
	},
	UserID:1,
	ShopID:1,
	ServicesID:1,
	AppointmentTime: &time.Time{},
	Shop:MockShop,
	User:MockUSer,


}

var MockSession =Session{
	Model:gorm.Model{
		ID:1,
		CreatedAt: time.Time{},
	},
	SessionId:"oenfebsdovdld dvn",
	UUID:1,
	Expires:10,
	SigningKey: []byte("sdddddddddddddddddddddddaskdnkjas"),

}

var MockReview = Review {
	Model:gorm.Model{
		ID:1,
		CreatedAt: time.Time{},
	},
	UserID:1,
	ShopID:1,
	Shop:MockShop,
	User:MockUSer,
	Review:"Mock Review 01",
	Reply: "Mock Reply 01",
	Rating:5,
}

