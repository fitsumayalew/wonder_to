package main

import (
	"net/http"
	"xCut/entity"
	"xCut/xcut/http/handler"

	shopRepoImport "xCut/shop/repository"
	shopServiceImport "xCut/shop/service"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/julienschmidt/httprouter"
)

func createTables(dbconn *gorm.DB) []error {
	errs := dbconn.DropTable(
		&entity.User{},
		&entity.Appointment{},
		&entity.Review{},
		&entity.Role{},
		&entity.Service{},
		&entity.Session{},
		&entity.Shop{},
	).GetErrors()

	errs = dbconn.CreateTable(
		&entity.User{},
		&entity.Appointment{},
		&entity.Review{},
		&entity.Role{},
		&entity.Service{},
		&entity.Session{},
		&entity.Shop{},
	).GetErrors()
	errs = dbconn.Create(&entity.Role{ID: 1, Name: "USER"}).GetErrors()
	errs = dbconn.Create(&entity.Role{ID: 2, Name: "ADMIN"}).GetErrors()
	if errs != nil {
		return errs
	}

	return nil
}
func main() {
	dbconn, err := gorm.Open("postgres", "postgres://postgres:Bangtan123@localhost/xcut?sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	createTables(dbconn)

	shopRepo2 := shopRepoImport.NewShopGormRepo(dbconn)
	shopService2 := shopServiceImport.NewShopService(shopRepo2)
	APIRouter := httprouter.New()
	APIRouter.ServeFiles("/assets/*filepath", http.Dir("../ui/assets"))
	ShopHandler := handler.NewShopHandler(shopService2)
	APIRouter.GET("/api/shops", ShopHandler.GetShops)
	APIRouter.GET("/api/shops/:id", ShopHandler.GetSingleshop)
	APIRouter.DELETE("/api/shops/:id", ShopHandler.DeleteShop)
	APIRouter.PUT("/api/shops/:id", ShopHandler.UpdateShop)
	APIRouter.POST("/api/shops", ShopHandler.PostShop)
	http.ListenAndServe(":8181", APIRouter)
}
