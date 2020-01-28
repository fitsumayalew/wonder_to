package appointment

import (
	"github.com/jinzhu/gorm"
	"time"
	"xCut/appointment"
	"xCut/entity"
)

type AppointmentGormRepo struct {
	conn *gorm.DB
}


// NewAppointmenttGormRepo returns new object of AppointmentGormRepo
func NewAppointmentGormRepo(db *gorm.DB) appointment.AppointmentRepository {
	return &AppointmentGormRepo{conn: db}
}


func (revRepo *AppointmentGormRepo) GetUpcomingShopAppointments(shopID uint) ([]entity.Appointment, []error) {
	appointments := []entity.Appointment{}
	errs := revRepo.conn.Find(&appointments, "shop_id=?", shopID).Where("AppointmentTime>?",time.Now()).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return appointments, errs
}

func (revRepo *AppointmentGormRepo) GetUserAppointments(userID uint) ([]entity.Appointment, []error) {
	appointments := []entity.Appointment{}
	errs := revRepo.conn.Find(&appointments, "user_id=?", userID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return appointments, errs
}


// GetAppointmentByShopID retrieve a Appointment from the database by its shopid
func (revRepo *AppointmentGormRepo) GetAppointments(shopID uint) ([]entity.Appointment, []error) {
	appointments := []entity.Appointment{}
	errs := revRepo.conn.Set("gorm:auto_preload", true).Find(&appointments, "shop_id=?", shopID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return appointments, errs
}

// UpdateAppointment updates a given Appointment in the database
func (revRepo *AppointmentGormRepo) UpdateAppointment(appointment *entity.Appointment) (*entity.Appointment, []error) {
	errs := revRepo.conn.Save(appointment).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}

func (revRepo *AppointmentGormRepo) GetAppointment(id uint) (*entity.Appointment, []error) {
	appointment := entity.Appointment{}
	errs := revRepo.conn.Set("gorm:auto_preload", true).First(&appointment, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &appointment, errs
}

// DeleteAppointment deletes a given Appointment from the database
func (revRepo *AppointmentGormRepo) DeleteAppointment(id uint) (*entity.Appointment, []error) {
	appointment, errs := revRepo.GetAppointment(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = revRepo.conn.Delete(appointment, appointment.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}

// StoreAppointment stores a given Appointment in the database
func (revRepo *AppointmentGormRepo) StoreAppointment(appointment *entity.Appointment) (*entity.Appointment, []error) {
	errs := revRepo.conn.Create(appointment).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}
