package appointment

import (
	"xCut/entity"
)

type AppointmentRepository interface {
	GetAppointments(shopID uint) ([]entity.Appointment, []error)
	GetAppointment(id uint) (*entity.Appointment, []error)
	GetUpcomingShopAppointments(shopID uint) ([]entity.Appointment, []error)
	GetUserAppointments(userID uint) ([]entity.Appointment, []error)
	StoreAppointment(appointment *entity.Appointment) (*entity.Appointment, []error)
	UpdateAppointment(appointment *entity.Appointment) (*entity.Appointment, []error)
	DeleteAppointment(id uint) (*entity.Appointment, []error)
}
