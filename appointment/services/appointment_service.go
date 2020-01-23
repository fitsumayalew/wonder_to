package appointment

import (
	"xCut/appointment"
	"xCut/entity"
)

type AppointmentService struct {
	appointmentRepo appointment.AppointmentService
}


func NewAppointmentService(appointmentRepository appointment.AppointmentRepository ) appointment.AppointmentRepository {
	return &AppointmentService{appointmentRepo: appointmentRepository}
}



func (appointmentService AppointmentService) DeleteAppointment(id uint) (*entity.Appointment, []error) {
	appointment, errs := appointmentService.appointmentRepo.DeleteAppointment(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}

func (appointmentService AppointmentService) GetAppointment(id uint) (*entity.Appointment, []error) {
	appointment, errs := appointmentService.appointmentRepo.GetAppointment(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}

func (appointmentService AppointmentService) GetAppointments(shopID uint) ([]entity.Appointment, []error) {
	appointment, errs := appointmentService.appointmentRepo.GetAppointments(shopID)
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}

func (appointmentService AppointmentService) GetUpcomingShopAppointments(shopID uint) ([]entity.Appointment, []error) {
	appointment, errs := appointmentService.appointmentRepo.GetUpcomingShopAppointments(shopID)
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}

func (appointmentService AppointmentService) GetUserAppointments(userID uint) ([]entity.Appointment, []error) {
	appointment, errs := appointmentService.appointmentRepo.GetUserAppointments(userID)
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}

func (appointmentService AppointmentService) StoreAppointment(appointment *entity.Appointment) (*entity.Appointment, []error) {
	appointment, errs := appointmentService.appointmentRepo.StoreAppointment(appointment)
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}

func (appointmentService AppointmentService) UpdateAppointment(appointment *entity.Appointment) (*entity.Appointment, []error) {
	appointment, errs := appointmentService.appointmentRepo.UpdateAppointment(appointment)
	if len(errs) > 0 {
		return nil, errs
	}
	return appointment, errs
}