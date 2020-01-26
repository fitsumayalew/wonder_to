package service

import (
	"xCut/entity"
)

type ServicesService interface {
	GetServices() ([]entity.Service, []error)
	GetService(id uint) (*entity.Service, []error)
	StoreService(service *entity.Service) (*entity.Service, []error)
	UpdateService(service *entity.Service) (*entity.Service, []error)
	DeleteService(id uint) (*entity.Service, []error)
	GetServiceByShopID(shopID uint) ([]entity.Service, []error)
}
