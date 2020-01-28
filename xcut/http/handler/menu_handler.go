package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"xCut/appointment"
	. "xCut/constants"
	"xCut/entity"
	"xCut/form"
	"xCut/review"
	"xCut/rtoken"
	"xCut/search"
	"xCut/service"
	"xCut/shop"
	"xCut/util"
)

type MenuHandler struct {
	tmpl               *template.Template
	shopService        shop.ShopService
	reviewService      review.ReviewService
	servicesService    service.ServicesService
	appointmentService appointment.AppointmentService
	searchService      search.SearchService
	csrfSignKey        []byte
}

func NewMenuHandler(t *template.Template,
	shopService shop.ShopService,
	reviewService review.ReviewService,
	servicesService service.ServicesService,
	appointmentService appointment.AppointmentService,
	searchService search.SearchService,
	csrfSignKey []byte,
) *MenuHandler {
	return &MenuHandler{tmpl: t, shopService: shopService, reviewService: reviewService,
		servicesService: servicesService, appointmentService: appointmentService, searchService: searchService,
		csrfSignKey: csrfSignKey}

}

func (menuHandler *MenuHandler) Index(w http.ResponseWriter, r *http.Request) {
	menuHandler.tmpl.ExecuteTemplate(w, "user.index.layout", menuHandler.csrfSignKey)
}

func (menuHandler *MenuHandler) Search(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		keyword := r.URL.Query().Get("address")
		var long float64
		var lat float64
		var err error
		var shops []entity.Shop
		if keyword == "" {
			long, err = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
			lat, err = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			shops, _ = menuHandler.searchService.GetByLocation(long, lat)

		} else {
			shops, _ = menuHandler.searchService.GetByName(keyword)
		}

		err = menuHandler.tmpl.ExecuteTemplate(w, "user.search.layout", shops)
		fmt.Println(err)
		return

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return

}

func (menuHandler *MenuHandler) BarberShop(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		shop, _ := menuHandler.shopService.GetShop(uint(id))
		reviews, _ := menuHandler.reviewService.GetReviewsByShopID(uint(id))
		services, _ := menuHandler.servicesService.GetServiceByShopID(uint(id))
		CSFRToken, err := rtoken.GenerateCSRFToken(menuHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		shopData := struct {
			Shop     entity.Shop
			Reviews  []entity.Review
			Services []entity.Service
			CSRF     string
		}{
			Reviews:  reviews,
			Shop:     *shop,
			Services: services,
			CSRF:     CSFRToken,
		}
		err = menuHandler.tmpl.ExecuteTemplate(w, "barbershop.layout", shopData)
		fmt.Println(err)
		return

	}

	if util.IsParsableFormPost(w, r, menuHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		reviewForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		reviewForm.ValidateRequiredFields(ReplyKey, RatingKey)
		rating, err := strconv.ParseUint(r.FormValue(RatingKey), 10, 32)
		shopID, err := strconv.ParseUint(r.FormValue(ReviewIDKey), 10, 32)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}
		currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)

		review := entity.Review{
			UserID: currentSession.UUID,
			ShopID: uint(shopID),
			Review: r.FormValue(ReplyKey),
			Reply:  "",
			Rating: uint(rating),
		}

		menuHandler.reviewService.StoreReview(&review)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return

}

func (menuHandler *MenuHandler) Review(w http.ResponseWriter, r *http.Request) {
	var shopID uint64
	if r.Method == http.MethodGet {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		shop, _ := menuHandler.shopService.GetShop(uint(id))
		CSFRToken, err := rtoken.GenerateCSRFToken(menuHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		shopData := struct {
			Shop    entity.Shop
			CSRF    string
			Values  url.Values
			VErrors form.ValidationErrors
		}{
			Shop:    *shop,
			CSRF:    CSFRToken,
			Values:  nil,
			VErrors: nil,
		}
		err = menuHandler.tmpl.ExecuteTemplate(w, "user.review.layout", shopData)
		fmt.Println(err)
		return

	}

	if util.IsParsableFormPost(w, r, menuHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		reviewForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		reviewForm.ValidateRequiredFields(ReviewKey, RatingKey)
		rating, err := strconv.ParseUint(r.FormValue(RatingKey), 10, 32)
		shopID, err = strconv.ParseUint(r.FormValue(ShopIDKey), 10, 32)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}
		currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)

		review := entity.Review{
			UserID: currentSession.UUID,
			ShopID: uint(shopID),
			Review: r.FormValue(ReviewKey),
			Reply:  "",
			Rating: uint(rating),
		}

		shop, _ := menuHandler.shopService.GetShop(uint(shopID))
		shop.Rating = (shop.Rating + uint(rating)) / 2
		menuHandler.shopService.UpdateShop(shop)
		menuHandler.reviewService.StoreReview(&review)

	}
	http.Redirect(w, r, "/barbershop?id="+fmt.Sprint(shopID), http.StatusSeeOther)
	return

}

func (menuHandler *MenuHandler) Appointment(w http.ResponseWriter, r *http.Request) {

	var shopID uint64
	if r.Method == http.MethodGet {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		shop, _ := menuHandler.shopService.GetShop(uint(id))
		services, _ := menuHandler.servicesService.GetServiceByShopID(uint(id))
		CSFRToken, err := rtoken.GenerateCSRFToken(menuHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		shopData := struct {
			Shop     entity.Shop
			CSRF     string
			Services []entity.Service
			Values   url.Values
			VErrors  form.ValidationErrors
		}{
			Shop:     *shop,
			CSRF:     CSFRToken,
			Services: services,
			Values:   nil,
			VErrors:  nil,
		}
		err = menuHandler.tmpl.ExecuteTemplate(w, "user.appointment.layout", shopData)
		fmt.Println(err)
		return

	}

	if util.IsParsableFormPost(w, r, menuHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		appointmentForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		appointmentForm .ValidateRequiredFields("mydate", "shopID","selectOne")
		shopID, err = strconv.ParseUint(r.FormValue("shopID"), 10, 32)
		serviceID, err := strconv.ParseUint(r.FormValue("selectOne"), 10, 32)
		//a:=strings.Split(r.FormValue("mydate"),"T")
		//b:=strings.Split(a[0],"-")
		//year:= b[0]
		//month:= b[1]
		//day:= b[2]
		//c:=strings.Split(a[1],":")
		//hour:= c[0]
		//minute:= c[1]

		mytime,err:= time.Parse(time.RFC3339,fmt.Sprint(r.FormValue("mydate"),":00Z"))
		if mytime.Before(time.Now()){
			appointmentForm.VErrors.Add("mydate","Can't input past date")
		}
		if !appointmentForm.IsValid() {
			menuHandler.tmpl.ExecuteTemplate(w, "user.appointment.layout", appointmentForm)
			return
		}

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}
		currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)

		appointment := entity.Appointment{

			UserID: currentSession.UUID,
			ShopID: uint(shopID),
			ServicesID: uint(serviceID),
			AppointmentTime: &mytime,
		}


		menuHandler.appointmentService.StoreAppointment(&appointment)

	}
	http.Redirect(w, r, "/barbershop?id="+fmt.Sprint(shopID), http.StatusSeeOther)
	return

}
