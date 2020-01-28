package appointment

import (
	"github.com/jinzhu/gorm"
	"errors"
	"xCut/appointment"
	"xCut/entity"
)

type MockAppointmentRepo struct {
	conn *gorm.DB
}

func (mappRepo *MockAppointmentRepo) GetAppointments(shopID uint) ([]entity.Appointment, []error) {
	Appo := []entity.Appointment{entity.MockAppointment}
	if shopID != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return Appo, nil
}

// NewAppointmenttGormRepo returns new object of MockAppointmentRepo
func NewMockAppointmentRepo(db *gorm.DB) appointment.AppointmentRepository {
	return &MockAppointmentRepo{conn: db}
}


func (mappRepo *MockAppointmentRepo) GetUpcomingShopAppointments(shopID uint) ([]entity.Appointment, []error) {
	Appo := []entity.Appointment{entity.MockAppointment}
	if shopID != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return Appo, nil
}

func (mappRepo *MockAppointmentRepo) GetUserAppointments(userID uint) ([]entity.Appointment, []error) {
	Appo := []entity.Appointment{entity.MockAppointment}
	if userID != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return Appo, nil
}


// GetAppointmentByShopID retrieve a Appointment from the database by its shopid
func (mappRepo *MockAppointmentRepo) GetAppointmentsByShopID(ShopID uint) ([]entity.Appointment, []error) {
	Appo := []entity.Appointment{entity.MockAppointment}

	if ShopID == 1 {
		return Appo, nil
	}
	return nil, []error{errors.New("Not found")}
}

// UpdateAppointment updates a given Appointment in the database
func (mappRepo *MockAppointmentRepo) UpdateAppointment(Appointment *entity.Appointment) (*entity.Appointment, []error) {
	Appo := entity.MockAppointment
	return &Appo, nil
}

func (mappRepo *MockAppointmentRepo) GetAppointment(id uint) (*entity.Appointment, []error) {
	Appo := entity.MockAppointment
	if id != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return &Appo, nil
}
// DeleteAppointment deletes a given Appointment from the database
func (mappRepo *MockAppointmentRepo) DeleteAppointment(id uint) (*entity.Appointment, []error) {
	Appo := entity.MockAppointment
	if id != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return &Appo, nil
}

// StoreAppointment stores a given Appointment in the database
func (mappRepo *MockAppointmentRepo) StoreAppointment(appointment *entity.Appointment) (*entity.Appointment, []error) {
	Appo := appointment
	return Appo,nil
}