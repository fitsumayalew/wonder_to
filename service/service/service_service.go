package service

import (
	"xCut/entity"
	"xCut/service"
)

type ServicesService struct {
	serviceRepo service.ServiceRepository
}

func (s ServicesService) GetServices() ([]entity.Service, []error) {
	services, errs := s.serviceRepo.GetServices()
	if len(errs) > 0 {
		return nil, errs
	}
	return services, errs
}

func (s ServicesService) GetService(id uint) (*entity.Service, []error) {
	services, errs := s.serviceRepo.GetService(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return services, errs
}

func (s ServicesService) StoreService(service *entity.Service) (*entity.Service, []error) {
	services, errs := s.serviceRepo.StoreService(service)
	if len(errs) > 0 {
		return nil, errs
	}
	return services, errs
}

func (s ServicesService) UpdateService(service *entity.Service) (*entity.Service, []error) {
	services, errs := s.serviceRepo.UpdateService(service)
	if len(errs) > 0 {
		return nil, errs
	}
	return services, errs
}

func (s ServicesService) DeleteService(id uint) (*entity.Service, []error) {
	services, errs := s.serviceRepo.DeleteService(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return services, errs
}

func (s ServicesService) GetServiceByShopID(shopID uint) ([]entity.Service, []error) {
	services, errs := s.serviceRepo.GetServiceByShopID(shopID)
	if len(errs) > 0 {
		return nil, errs
	}
	return services, errs
}

func NewServiceService(serviceRepo service.ServiceRepository) service.ServicesService {
	return &ServicesService{serviceRepo: serviceRepo}
}
