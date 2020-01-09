package handler

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strconv"
	"xcut/entity"
	"xcut/form"
	"xcut/rtoken"
	"xcut/shop"
	"xcut/util"
)

type AdminDashboardHandler struct {
	tmpl        *template.Template
	shopService shop.ShopService
	csrfSignKey []byte
}





func NewAdminDashboardHandler(
	t *template.Template,
	shopService shop.ShopService,
	csrfSignKey []byte,
) *AdminDashboardHandler {
	return &AdminDashboardHandler{tmpl: t, shopService: shopService, csrfSignKey: csrfSignKey}
}

func (adminDashboardHandler *AdminDashboardHandler) AdminIndex(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop,errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0{
		http.Redirect(w, r, "/finishSignup", http.StatusSeeOther)
		return
	}


	err := adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.index.layout", shop)
	fmt.Println(err)
}

func (adminDashboardHandler *AdminDashboardHandler) AdminBasicInfo(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	shop,errs := adminDashboardHandler.shopService.GetShopByUserID(currentSession.UUID)
	if len(errs) > 0{
		http.Redirect(w, r, "/finishSignup", http.StatusSeeOther)
		return
	}


	err := adminDashboardHandler.tmpl.ExecuteTemplate(w, "admin.basic.layout", shop)
	fmt.Println(err)
}


func (adminDashboardHandler *AdminDashboardHandler) AdminSignUp(w http.ResponseWriter,r *http.Request){
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	//If it's requesting the login page return CSFR Signed token with the form
	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(adminDashboardHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		adminDashboardHandler.tmpl.ExecuteTemplate(w, "signup.shop.layout",form.Input{CSRF: CSFRToken})
		return
	}
	//Only reply to forms that have that are parsable and have valid csfrToken
	if adminDashboardHandler.isParsableFormPost(w, r) {

		//Validate form data
		signUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		signUpForm.ValidateRequiredFields(shopNameKey,phoneKey,cityKey,addressKey,latKey,longKey)
		signUpForm.MatchesPattern(phoneKey,form.PhoneRX)
		//phone := r.FormValue(phoneKey)
		//user, errs := adminDashboardHandler.shopService.

		if !signUpForm.IsValid() {
			adminDashboardHandler.tmpl.ExecuteTemplate(w, "signup.layout", signUpForm)
			return
		}

		long,_ := strconv.ParseFloat(r.FormValue(longKey),32)
		lat,_ := strconv.ParseFloat(r.FormValue(latKey),32)
		shop := entity.Shop{
			Model:    gorm.Model{},
			Name:     r.FormValue(shopNameKey),
			City:     r.FormValue(cityKey),
			Lat:      float32(lat),
			Long:     float32(long),
			Address:  r.FormValue(addressKey),
			Phone:    r.FormValue(phoneKey),
			Website:  "",
			Image:    "",
			UserID:   currentSession.UUID,
			Services: nil,
		}
		adminDashboardHandler.shopService.StoreShop(&shop)

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}


func (adminDashboardHandler *AdminDashboardHandler) isParsableFormPost(w http.ResponseWriter, r *http.Request) bool {
	return r.Method == http.MethodPost &&
		util.ParseForm(w, r) &&
		rtoken.IsCSRFValid(r.FormValue(csrfKey), adminDashboardHandler.csrfSignKey)
}