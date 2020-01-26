package service

import (
	"xCut/entity"
)

type ServiceRepository interface {
	GetServices() ([]entity.Service, []error)
	GetService(id uint) (*entity.Service, []error)
	StoreService(service *entity.Service) (*entity.Service, []error)
	UpdateService(service *entity.Service) (*entity.Service, []error)
	DeleteService(id uint) (*entity.Service, []error)
	GetServiceByShopID(ShopID uint) ([]entity.Service, []error)
}
