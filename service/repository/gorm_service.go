package repository

import (
	"github.com/jinzhu/gorm"
	"xCut/entity"
	"xCut/service"
)

type ServiceGormRepo struct {
	conn *gorm.DB
}

// NewReviewtGormRepo returns new object of ReviewGormRepo
func NewServiceGormRepo(db *gorm.DB) service.ServiceRepository {
	return &ServiceGormRepo{conn: db}
}

func (serviceRepo ServiceGormRepo) GetServices() ([]entity.Service, []error) {
	services := []entity.Service{}
	errs := serviceRepo.conn.Find(&services).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return services, errs
}

func (serviceRepo ServiceGormRepo) GetService(id uint) (*entity.Service, []error) {
	service := entity.Service{}
	errs := serviceRepo.conn.Set("gorm:auto_preload", true).First(&service, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &service, errs
}

func (serviceRepo ServiceGormRepo) StoreService(service *entity.Service) (*entity.Service, []error) {
	errs := serviceRepo.conn.Create(service).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return service, errs
}

func (serviceRepo ServiceGormRepo) UpdateService(service *entity.Service) (*entity.Service, []error) {
	errs := serviceRepo.conn.Save(service).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return service, errs
}

func (serviceRepo ServiceGormRepo) DeleteService(id uint) (*entity.Service, []error) {
	service, errs := serviceRepo.GetService(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = serviceRepo.conn.Delete(service, service.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return service, errs
}

func (serviceRepo ServiceGormRepo) GetServiceByShopID(ShopID uint) ([]entity.Service, []error) {
	services := []entity.Service{}
	errs := serviceRepo.conn.Set("gorm:auto_preload", true).Find(&services, "shop_id=?", ShopID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return services, errs
}


