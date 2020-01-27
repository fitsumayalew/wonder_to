package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	"net/http"
	appointmentRepoImport "xCut/appointment/repository"
	appointmentServiceImport "xCut/appointment/services"
	"xCut/entity"
	reviewRepoImport "xCut/review/repository"
	reviewServiceImport "xCut/review/service"
	"xCut/rtoken"
	searchRepoImport "xCut/search/repository"
	searchServiceImport "xCut/search/service"
	serviceRepoImport "xCut/service/repository"
	serviceServiceImport "xCut/service/service"
	shopRepoImport "xCut/shop/repository"
	shopServiceImport "xCut/shop/service"
	userRepoImport "xCut/user/repository"
	userServiceImport "xCut/user/service"
	"xCut/xcut/http/handler"
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
	fm := template.FuncMap{"processDate": func(time uint) string {
		hour := time / 60
		postfix := "AM"
		if hour > 12 {
			hour = hour - 12
			postfix = "PM"
		}
		minute := time % 60
		return fmt.Sprintf("%02d:%02d %s", hour, minute, postfix)

	}}

	// Create a template, add the function map, and parse the text.
	tmpl := template.Must(template.New("main").Funcs(fm).ParseGlob("ui/templates/*"))
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

	reviewRepo := reviewRepoImport.NewReviewGormRepo(dbconn)
	reviewService := reviewServiceImport.NewReviewService(reviewRepo)

	serviceRepo := serviceRepoImport.NewServiceGormRepo(dbconn)
	servicesService := serviceServiceImport.NewServiceService(serviceRepo)

	searchRepo := searchRepoImport.NewSearchGormRepo(dbconn)
	searchService := searchServiceImport.NewSearchService(searchRepo)

	appointmentRepo := appointmentRepoImport.NewAppointmentGormRepo(dbconn)
	appointmentService := appointmentServiceImport.NewAppointmentService(appointmentRepo)



	userHandler := handler.NewUserHandler(tmpl, userService, sessionService, roleService, csrfSignKey)
	adminDashboardHandler := handler.NewAdminDashboardHandler(tmpl, shopService,reviewService, servicesService,appointmentService,csrfSignKey)
	menuHandler := handler.NewMenuHandler(tmpl, shopService,reviewService, servicesService,appointmentService,searchService, csrfSignKey)

	fs := http.FileServer(http.Dir("ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.Handle("/admin", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminIndex))))
	http.HandleFunc("/login", userHandler.Login)
	http.Handle("/admin/finishSignup", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminSignUp))))
	http.Handle("/admin/basicInfo", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminBasicInfo))))
	http.Handle("/admin/basicInfoEdit", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminBasicInfoEdit))))
	http.Handle("/admin/services", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminServices))))
	http.Handle("/admin/appointments", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminAppointments))))
	http.Handle("/admin/reviews", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminReviews))))
	http.Handle("/admin/reply", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminReply))))


	http.HandleFunc("/", menuHandler.Index)
	http.HandleFunc("/search", menuHandler.Search)

	http.HandleFunc("/barbershop", menuHandler.BarberShop)


	http.Handle("/admin/services/new", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminServicesAdd))))
	http.Handle("/admin/services/update", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminServicesUpdate))))
	http.Handle("/admin/services/delete", userHandler.Authenticated(userHandler.Authorized(http.HandlerFunc(adminDashboardHandler.AdminServicesDelete))))


	http.Handle("/logout", userHandler.Authenticated(http.HandlerFunc(userHandler.Logout)))
	http.HandleFunc("/signup", userHandler.SignUp)
	http.ListenAndServe(":8181", nil)
}
