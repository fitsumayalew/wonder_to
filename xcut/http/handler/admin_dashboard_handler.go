package handler

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"
	. "xCut/constants"
	"xCut/entity"
	"xCut/form"
	"xCut/review"
	"xCut/rtoken"
	"xCut/service"
	"xCut/shop"
	"xCut/util"
)

type AdminDashboardHandler struct {
	tmpl            *template.Template
	shopService     shop.ShopService
	reviewService   review.ReviewService
	servicesService service.ServicesService
	csrfSignKey     []byte
}

func NewAdminDashboardHandler(
	t *template.Template,
	shopService shop.ShopService,
	reviewService review.ReviewService,
	servicesService service.ServicesService,
	csrfSignKey []byte,
) *AdminDashboardHandler {
	return &AdminDashboardHandler{tmpl: t, shopService: shopService, reviewService: reviewService, servicesService: servicesService, csrfSignKey: csrfSignKey}
}

func (adminDashboardHandler *AdminDashboardHandler) AdminIndex(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	err := adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.index.layout", shop)
	fmt.Println(err)
}

func (adminDashboardHandler *AdminDashboardHandler) AdminBasicInfo(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	err := adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.basic.layout", shop)
	fmt.Println(err)
}

func (adminDashboardHandler *AdminDashboardHandler) AdminBasicInfoEdit(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		shopMap := url.Values{}
		shopMap[ShopNameKey] = []string{shop.Name}
		shopMap[PhoneKey] = []string{shop.Phone}
		shopMap[CityKey] = []string{shop.City}
		shopMap[AddressKey] = []string{shop.Address}
		shopMap[WebsiteKey] = []string{shop.Website}
		shopMap[LatKey] = []string{fmt.Sprintf("%f", shop.Lat)}
		shopMap[LngKey] = []string{fmt.Sprintf("%f", shop.Long)}
		shopMap[CsrfKey] = []string{CSFRToken}
		shopMap[WeekdaysOpenHoursStart] = []string{fmt.Sprintf("%02d:%02d", shop.WeekDayOpenHour/60, shop.WeekDayOpenHour%60)}
		shopMap[WeekdaysOpenHoursEnd] = []string{fmt.Sprintf("%02d:%02d", shop.WeekDayCloseHour/60, shop.WeekDayCloseHour%60)}
		shopMap[WeekendsOpenHoursStart] = []string{fmt.Sprintf("%02d:%02d", shop.WeekendOpenHour/60, shop.WeekendOpenHour%60)}
		shopMap[WeekendsOpenHoursEnd] = []string{fmt.Sprintf("%02d:%02d", shop.WeekendCloseHour/60, shop.WeekendCloseHour%60)}
		signUpForm := form.Input{Values: shopMap, VErrors: form.ValidationErrors{}}

		adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.basic.edit.layout", signUpForm)
	}

	if util.IsParsableFormPost(w, r, adminDashboardHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		basicInfoEditForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		basicInfoEditForm.ValidateRequiredFields(ShopNameKey, PhoneKey, CityKey, AddressKey, LatKey, LngKey, WeekdaysOpenHoursEnd, WeekdaysOpenHoursStart, WeekendsOpenHoursEnd, WeekdaysOpenHoursStart)
		basicInfoEditForm.MatchesPattern(PhoneKey, form.PhoneRX)
		basicInfoEditForm.MatchesPattern(WebsiteKey, form.WebsiteRX)
		basicInfoEditForm.ValidateStartAndEnd(WeekdaysOpenHoursStart, WeekdaysOpenHoursEnd)
		basicInfoEditForm.ValidateStartAndEnd(WeekendsOpenHoursStart, WeekendsOpenHoursEnd)
		basicInfoEditForm.MatchesPattern(LatKey, form.LngLatRX)
		basicInfoEditForm.MatchesPattern(LngKey, form.LngLatRX)

		if !basicInfoEditForm.IsValid() {
			adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.basic.edit.layout", basicInfoEditForm)
			return
		}

		weekdaysOpenHoursStart, _ := time.Parse("15:06", r.FormValue(WeekdaysOpenHoursStart))
		weekdaysOpenHoursEnd, _ := time.Parse("15:06", r.FormValue(WeekdaysOpenHoursEnd))

		weekendOpenHoursStart, _ := time.Parse("15:06", r.FormValue(WeekendsOpenHoursStart))
		weekendOpenHoursEnd, _ := time.Parse("15:06", r.FormValue(WeekendsOpenHoursEnd))

		long, _ := strconv.ParseFloat(r.FormValue(LngKey), 32)
		lat, _ := strconv.ParseFloat(r.FormValue(LatKey), 32)

		if r.FormValue(ImageFile) != "" {
			mf, _, err := r.FormFile(ImageFile)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			fileName, err := util.GenerateFileName(&mf, r.FormValue(ImageFile))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			shop.Image = fileName
			if err == nil {
				shop.Image = fileName
				err = util.WriteFile(&mf, shop.Image)
			}
			if mf != nil {
				defer mf.Close()
			}
		}
		shop.Name = r.FormValue(ShopNameKey)
		shop.City = r.FormValue(CityKey)
		shop.Lat = lat
		shop.Long = long
		shop.Address = r.FormValue(AddressKey)
		shop.Phone = r.FormValue(PhoneKey)
		shop.Website = r.FormValue(WebsiteKey)
		shop.WeekDayOpenHour = uint(weekdaysOpenHoursStart.Hour()*60 + weekdaysOpenHoursStart.Minute())
		shop.WeekDayCloseHour = uint(weekdaysOpenHoursEnd.Hour()*60 + weekdaysOpenHoursEnd.Minute())
		shop.WeekendOpenHour = uint(weekendOpenHoursStart.Hour()*60 + weekendOpenHoursStart.Minute())
		shop.WeekendCloseHour = uint(weekendOpenHoursEnd.Hour()*60 + weekendOpenHoursStart.Minute())
		//shop.Services = nil

		// Save the user to the database
		_, ers := adminDashboardHandler.shopService.UpdateShop(shop)
		if len(ers) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/basicInfo", http.StatusSeeOther)

	}
}

func (adminDashboardHandler *AdminDashboardHandler) AdminSignUp(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	//If it's requesting the login page return CSFR Signed token with the form
	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		adminDashboardHandler.tmpl.ExecuteTemplate(w, "signup.shop.layout", form.Input{CSRF: CSFRToken})
		return
	}
	//Only reply to forms that have that are parsable and have valid csfrToken
	if util.IsParsableFormPost(w, r, adminDashboardHandler.csrfSignKey) {

		//Validate form data
		signUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		signUpForm.ValidateRequiredFields(ShopNameKey, PhoneKey, CityKey, AddressKey, LatKey, LngKey)
		signUpForm.MatchesPattern(PhoneKey, form.PhoneRX)
		//phone := r.FormValue(phoneKey)
		//user, errs := adminDashboardHandler.shopService.

		if !signUpForm.IsValid() {
			adminDashboardHandler.tmpl.ExecuteTemplate(w, "signup.layout", signUpForm)
			return
		}

		long, _ := strconv.ParseFloat(r.FormValue(LngKey), 32)
		lat, _ := strconv.ParseFloat(r.FormValue(LatKey), 32)
		shop := entity.Shop{
			Model:   gorm.Model{},
			Name:    r.FormValue(ShopNameKey),
			City:    r.FormValue(CityKey),
			Lat:     lat,
			Long:    long,
			Address: r.FormValue(AddressKey),
			Phone:   r.FormValue(PhoneKey),
			Website: "",
			Image:   "",
			UserID:  currentSession.UUID,
			//Services: nil,
		}
		adminDashboardHandler.shopService.StoreShop(&shop)

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func (adminDashboardHandler *AdminDashboardHandler) AdminAppointments(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		shopMap := url.Values{}
		shopMap[ShopNameKey] = []string{shop.Name}
		shopMap[PhoneKey] = []string{shop.Phone}
		shopMap[CityKey] = []string{shop.City}
		shopMap[AddressKey] = []string{shop.Address}
		shopMap[WebsiteKey] = []string{shop.Website}
		shopMap[LatKey] = []string{fmt.Sprintf("%f", shop.Lat)}
		shopMap[LngKey] = []string{fmt.Sprintf("%f", shop.Long)}
		shopMap[CsrfKey] = []string{CSFRToken}
		shopMap[WeekdaysOpenHoursStart] = []string{fmt.Sprintf("%02d:%02d", shop.WeekDayOpenHour/60, shop.WeekDayOpenHour%60)}
		shopMap[WeekdaysOpenHoursEnd] = []string{fmt.Sprintf("%02d:%02d", shop.WeekDayCloseHour/60, shop.WeekDayCloseHour%60)}
		shopMap[WeekendsOpenHoursStart] = []string{fmt.Sprintf("%02d:%02d", shop.WeekendOpenHour/60, shop.WeekendOpenHour%60)}
		shopMap[WeekendsOpenHoursEnd] = []string{fmt.Sprintf("%02d:%02d", shop.WeekendCloseHour/60, shop.WeekendCloseHour%60)}
		signUpForm := form.Input{Values: shopMap, VErrors: form.ValidationErrors{}}

		adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.basic.edit.layout", signUpForm)
	}

	if util.IsParsableFormPost(w, r, adminDashboardHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		basicInfoEditForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		basicInfoEditForm.ValidateRequiredFields(ShopNameKey, PhoneKey, CityKey, AddressKey, LatKey, LngKey, WeekdaysOpenHoursEnd, WeekdaysOpenHoursStart, WeekendsOpenHoursEnd, WeekdaysOpenHoursStart)
		basicInfoEditForm.MatchesPattern(PhoneKey, form.PhoneRX)
		basicInfoEditForm.MatchesPattern(WebsiteKey, form.WebsiteRX)
		basicInfoEditForm.ValidateStartAndEnd(WeekdaysOpenHoursStart, WeekdaysOpenHoursEnd)
		basicInfoEditForm.ValidateStartAndEnd(WeekendsOpenHoursStart, WeekendsOpenHoursEnd)
		basicInfoEditForm.MatchesPattern(LatKey, form.LngLatRX)
		basicInfoEditForm.MatchesPattern(LngKey, form.LngLatRX)

		if !basicInfoEditForm.IsValid() {
			adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.basic.edit.layout", basicInfoEditForm)
			return
		}

		weekdaysOpenHoursStart, _ := time.Parse("15:06", r.FormValue(WeekdaysOpenHoursStart))
		weekdaysOpenHoursEnd, _ := time.Parse("15:06", r.FormValue(WeekdaysOpenHoursEnd))

		weekendOpenHoursStart, _ := time.Parse("15:06", r.FormValue(WeekendsOpenHoursStart))
		weekendOpenHoursEnd, _ := time.Parse("15:06", r.FormValue(WeekendsOpenHoursEnd))

		long, _ := strconv.ParseFloat(r.FormValue(LngKey), 32)
		lat, _ := strconv.ParseFloat(r.FormValue(LatKey), 32)

		if r.FormValue(ImageFile) != "" {
			mf, _, err := r.FormFile(ImageFile)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			fileName, err := util.GenerateFileName(&mf, r.FormValue(ImageFile))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			shop.Image = fileName
			if err == nil {
				shop.Image = fileName
				err = util.WriteFile(&mf, shop.Image)
			}
			if mf != nil {
				defer mf.Close()
			}
		}
		shop.Name = r.FormValue(ShopNameKey)
		shop.City = r.FormValue(CityKey)
		shop.Lat = lat
		shop.Long = long
		shop.Address = r.FormValue(AddressKey)
		shop.Phone = r.FormValue(PhoneKey)
		shop.Website = r.FormValue(WebsiteKey)
		shop.WeekDayOpenHour = uint(weekdaysOpenHoursStart.Hour()*60 + weekdaysOpenHoursStart.Minute())
		shop.WeekDayCloseHour = uint(weekdaysOpenHoursEnd.Hour()*60 + weekdaysOpenHoursEnd.Minute())
		shop.WeekendOpenHour = uint(weekendOpenHoursStart.Hour()*60 + weekendOpenHoursStart.Minute())
		shop.WeekendCloseHour = uint(weekendOpenHoursEnd.Hour()*60 + weekendOpenHoursStart.Minute())
		//shop.Services = nil

		// Save the user to the database
		_, ers := adminDashboardHandler.shopService.UpdateShop(shop)
		if len(ers) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/basicInfo", http.StatusSeeOther)

	}
}
func (adminDashboardHandler *AdminDashboardHandler) AdminReviews(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		reviews, errs := adminDashboardHandler.reviewService.GetReviewsByShopID(shop.ID)

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		reviewsData := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			Reviews []entity.Review
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			Reviews: reviews,
			CSRF:    CSFRToken,
		}
		adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.reviews.layout", reviewsData)
	}
}

func (adminDashboardHandler *AdminDashboardHandler) AdminReply(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		id, err := strconv.Atoi(r.URL.Query().Get("id"))

		if err != nil {
			http.Redirect(w, r, "/admin/reviews", http.StatusSeeOther)
			return
		}

		review, _ := adminDashboardHandler.reviewService.GetReview(uint(id))

		replyData := struct {
			Values   url.Values
			VErrors  form.ValidationErrors
			ReviewID int
			Review   entity.Review
			CSRF     string
		}{
			Values:   nil,
			VErrors:  nil,
			ReviewID: id,
			Review:   *review,
			CSRF:     CSFRToken,
		}
		adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.reply.layout", replyData)
	}

	if util.IsParsableFormPost(w, r, adminDashboardHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		replyForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		replyForm.ValidateRequiredFields(ReplyKey, ReviewIDKey)
		reviewID, err := strconv.ParseUint(r.FormValue(ReviewIDKey), 10, 32)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}

		review, _ := adminDashboardHandler.reviewService.GetReview(uint(reviewID))
		if shop.ID != review.ShopID {
			replyForm.VErrors.Add(ReplyKey, "Review not found")
		}

		if !replyForm.IsValid() {
			adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.reply.layout", replyForm)
			return
		}

		review.Reply = r.FormValue(ReplyKey)
		_, ers := adminDashboardHandler.reviewService.UpdateReview(review)
		if len(ers) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/reviews", http.StatusSeeOther)

	}
}

func (adminDashboardHandler *AdminDashboardHandler) AdminServices(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		services, _ := adminDashboardHandler.servicesService.GetServiceByShopID(shop.ID)

		adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.service.layout", services)
	}

}

func (adminDashboardHandler *AdminDashboardHandler) AdminServicesAdd(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}
	fmt.Println(shop.Name)

	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		servicesData := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    CSFRToken,
		}

		adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.service.new.layout", servicesData)
		return
	}

	if util.IsParsableFormPost(w, r, adminDashboardHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		addServiceForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		addServiceForm.ValidateRequiredFields(ServiceETCKey, ServiceNameKey, ServicePriceKey)
		addServiceForm.MatchesPattern(PhoneKey, form.PhoneRX)
		addServiceForm.MatchesPattern(ServicePriceKey, form.LngLatRX)
		price, _ := strconv.ParseFloat(r.FormValue(ServicePriceKey), 32)
		etc, _ := strconv.Atoi(r.FormValue((ServiceETCKey)))
		if !addServiceForm.IsValid() {
			adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.service.new.layout", addServiceForm)
			return
		}

		fileName := ""
		mf, fh, err := r.FormFile(ImageFile)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fileName, err = util.GenerateFileName(&mf, fh.Filename)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		shop.Image = fileName
		if err == nil {
			shop.Image = fileName
			err = util.WriteFile(&mf, shop.Image)
		}
		if mf != nil {
			defer mf.Close()
		}

		service := entity.Service{
			ShopID:        shop.ID,
			Name:          r.FormValue(ServiceNameKey),
			Price:         float32(price),
			EstimatedTime: uint(etc),
			Image:         fileName,
			Shop:          *shop,
		}

		_, ers := adminDashboardHandler.servicesService.StoreService(&service)
		if len(ers) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/services/new", http.StatusSeeOther)

	}
}

func (adminDashboardHandler *AdminDashboardHandler) AdminServicesUpdate(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
		return
	}

	serviceObj, errs := adminDashboardHandler.servicesService.GetService(uint(id))

	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		s := url.Values{}
		s.Add(ServiceNameKey, serviceObj.Name)
		s.Add(ServicePriceKey, fmt.Sprint(serviceObj.Price))
		s.Add(ServiceETCKey, fmt.Sprint(serviceObj.EstimatedTime))

		servicesData := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  s,
			VErrors: nil,
			CSRF:    CSFRToken,
		}

		adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.service.new.layout", servicesData)
		return
	}

	if util.IsParsableFormPost(w, r, adminDashboardHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		serviceUpdateForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		addServiceForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		addServiceForm.ValidateRequiredFields(ServiceETCKey, ServiceNameKey, ServicePriceKey)
		addServiceForm.MatchesPattern(PhoneKey, form.PhoneRX)
		addServiceForm.MatchesPattern(ServicePriceKey, form.LngLatRX)
		price, _ := strconv.ParseFloat(r.FormValue(ServicePriceKey), 32)
		etc, _ := strconv.Atoi(r.FormValue((ServiceETCKey)))
		if !serviceUpdateForm.IsValid() {
			adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.basic.edit.layout", serviceUpdateForm)
			return
		}

		mf, fm, err := r.FormFile(ImageFile)
		if mf != nil {
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			fileName, err := util.GenerateFileName(&mf, fm.Filename)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			shop.Image = fileName
			if err == nil {
				shop.Image = fileName
				err = util.WriteFile(&mf, shop.Image)
			}
			if mf != nil {
				defer mf.Close()
			}
			serviceObj.Image = fileName
		}
		serviceObj.Name = r.FormValue(ServiceNameKey)
		serviceObj.EstimatedTime = uint(etc)
		serviceObj.Price = float32(price)
		// Save the user to the database
		_, ers := adminDashboardHandler.servicesService.UpdateService(serviceObj)
		if len(ers) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/services", http.StatusSeeOther)

	}
}

func (adminDashboardHandler *AdminDashboardHandler) AdminServicesDelete(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop, errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0 {
		http.Redirect(w, r, "/admin/finishSignup", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		servicesData := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			//Service entity.Service
			CSRF string
		}{
			Values:  nil,
			VErrors: nil,
			//Service :  nil,
			CSRF: CSFRToken,
		}

		adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.services.layout", servicesData)
	}

	if util.IsParsableFormPost(w, r, adminDashboardHandler.csrfSignKey) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		basicInfoEditForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		basicInfoEditForm.ValidateRequiredFields(ShopNameKey, PhoneKey, CityKey, AddressKey, LatKey, LngKey, WeekdaysOpenHoursEnd, WeekdaysOpenHoursStart, WeekendsOpenHoursEnd, WeekdaysOpenHoursStart)
		basicInfoEditForm.MatchesPattern(PhoneKey, form.PhoneRX)
		basicInfoEditForm.MatchesPattern(WebsiteKey, form.WebsiteRX)
		basicInfoEditForm.ValidateStartAndEnd(WeekdaysOpenHoursStart, WeekdaysOpenHoursEnd)
		basicInfoEditForm.ValidateStartAndEnd(WeekendsOpenHoursStart, WeekendsOpenHoursEnd)
		basicInfoEditForm.MatchesPattern(LatKey, form.LngLatRX)
		basicInfoEditForm.MatchesPattern(LngKey, form.LngLatRX)

		if !basicInfoEditForm.IsValid() {
			adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.basic.edit.layout", basicInfoEditForm)
			return
		}

		weekdaysOpenHoursStart, _ := time.Parse("15:06", r.FormValue(WeekdaysOpenHoursStart))
		weekdaysOpenHoursEnd, _ := time.Parse("15:06", r.FormValue(WeekdaysOpenHoursEnd))

		weekendOpenHoursStart, _ := time.Parse("15:06", r.FormValue(WeekendsOpenHoursStart))
		weekendOpenHoursEnd, _ := time.Parse("15:06", r.FormValue(WeekendsOpenHoursEnd))

		long, _ := strconv.ParseFloat(r.FormValue(LngKey), 32)
		lat, _ := strconv.ParseFloat(r.FormValue(LatKey), 32)

		if r.FormValue(ImageFile) != "" {
			mf, _, err := r.FormFile(ImageFile)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			fileName, err := util.GenerateFileName(&mf, r.FormValue(ImageFile))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			shop.Image = fileName
			if err == nil {
				shop.Image = fileName
				err = util.WriteFile(&mf, shop.Image)
			}
			if mf != nil {
				defer mf.Close()
			}
		}
		shop.Name = r.FormValue(ShopNameKey)
		shop.City = r.FormValue(CityKey)
		shop.Lat = lat
		shop.Long = long
		shop.Address = r.FormValue(AddressKey)
		shop.Phone = r.FormValue(PhoneKey)
		shop.Website = r.FormValue(WebsiteKey)
		shop.WeekDayOpenHour = uint(weekdaysOpenHoursStart.Hour()*60 + weekdaysOpenHoursStart.Minute())
		shop.WeekDayCloseHour = uint(weekdaysOpenHoursEnd.Hour()*60 + weekdaysOpenHoursEnd.Minute())
		shop.WeekendOpenHour = uint(weekendOpenHoursStart.Hour()*60 + weekendOpenHoursStart.Minute())
		shop.WeekendCloseHour = uint(weekendOpenHoursEnd.Hour()*60 + weekendOpenHoursStart.Minute())
		//shop.Services = nil

		// Save the user to the database
		_, ers := adminDashboardHandler.shopService.UpdateShop(shop)
		if len(ers) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/basicInfo", http.StatusSeeOther)

	}
}
