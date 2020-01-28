package entity
 import (

"time"
"github.com/fitsumayalew/xCut/entity"
)

type MockRole = Role{
	ID:  1,
	Name:  "Mock Role 01",
}

type MockService = Service{
  	ID:1,
	Name:"Mock Service",
	ShopID:1,
	Shop:mock_shop,
	Price:75,
	Image:"mock_services.png",
	EstimatedTime: 10,
	CreatedAt: time.Time{},
}

type MockUSer = User{
	ID:1,
	FullName:"Mock User 01",
	Email:"mockuser@gmail.com",
	Password:"P@$$w0rd",
	Phone:"0901010101",
	NumberOfReviews:3,
	Image:"mock_user.png",
	RoleID: 1,
	CreatedAt: time.Time{},

	}

//Shop Represents the the shops the user views
type MockShop = Shop{
	ID:1,
	Name:" Mock Shop 01 ",
	City:"Mock City 01",
	Lat:10.2,
	Long:10.5,
	Address:"Mock Adress 01"
	Phone:"0901010101",
	Website:"WWW.mocksite.com",
	Image:"mock_item.png",
	UserID:1,
	WeekDayOpenHour:2,
	WeekDayCloseHour:3,
	WeekendOpenHour:2,
	WeekendCloseHour:3,
	services: []services{},
	Reviews    []Review{},
	Appointments  []Appointment{},
	CreatedAt: time.Time{},
}


type MockAppointment =Appointment{
	ID:1,
	UserID:1,
	ShopID:1,
	ServicesID:1,
	CreatedAt: time.Time{},
	AppointmentTime: time.Time{},
	Shop:mock_shop,
	User:mock_user,


}

type MockSession =Session{
	ID: 1,
	SessionId:"oenfebsdovdld dvn",
	UUID:1,
	Expires:10,
	SigningKey: []byte "MockSignUpkey",
	CreatedAt: time.Time{},

}

type MockReview = Review {
	ID:1,
	UserID:1,
	ShopID:1,
	Shop:mock_shop,
	User:mock_user,
	Review:"Mock Review 01",
	Reply: "Mock Reply 01",
	Rating:5,
	CreatedAt: time.Time{},
}

