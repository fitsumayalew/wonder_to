package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	"net/http"
	"xcut/entity"
	"xcut/rtoken"
	shopRepoImport "xcut/shop/repository"
	shopServiceImport "xcut/shop/service"
	userRepoImport "xcut/user/repository"
	userServiceImport "xcut/user/service"
	"xcut/xcut/http/handler"
)

func createTables(dbconn *gorm.DB) []error {
	errs := dbconn.DropTable(
		&entity.User{},
		&entity.Appointments{},
		&entity.Reply{},
		&entity.Review{},
		&entity.Role{},
		&entity.Services{},
		&entity.Session{},
		&entity.Shop{},
	).GetErrors()

	errs = dbconn.CreateTable(
		&entity.User{},
		&entity.Appointments{},
		&entity.Reply{},
		&entity.Review{},
		&entity.Role{},
		&entity.Services{},
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
	tmpl := template.Must(template.ParseGlob("ui/templates/*"))
	dbconn, err := gorm.Open("postgres", "postgres://postgres:root@localhost/xcut?sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	//createTables(dbconn)

	csrfSignKey := []byte(rtoken.GenerateRandomID(32))
	userRepo := userRepoImport.NewUserGormRepo(dbconn)
	userService := userServiceImport.NewUserService(userRepo)

	sessionRepo := userRepoImport.NewSessionGormRepo(dbconn)
	sessionService := userServiceImport.NewSessionService(sessionRepo)

	roleRepo := userRepoImport.NewRoleGormRepo(dbconn)
	roleService := userServiceImport.NewRoleService(roleRepo)

	shopRepo := shopRepoImport.NewShopGormRepo(dbconn)
	shopService := shopServiceImport.NewShopService(shopRepo)

	userHandler := handler.NewUserHandler(tmpl, userService, sessionService, roleService, csrfSignKey)
	adminDashboardHandler := handler.NewAdminDashboardHandler(tmpl,shopService,csrfSignKey)

	fs := http.FileServer(http.Dir("ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", userHandler.Index)
	http.Handle("/admin", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminIndex))))
	http.HandleFunc("/login", userHandler.Login)
	http.Handle("/finishSignup",userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminSignUp))))
	http.Handle("/basicInfo",userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminBasicInfo))))
	http.Handle("/logout", userHandler.Authenticated(http.HandlerFunc(userHandler.Logout)))
	http.HandleFunc("/signup", userHandler.SignUp)
	http.ListenAndServe(":8181", nil)
}
