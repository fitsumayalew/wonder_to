package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	shop "xCut/shop/repository"
	shopServiceImport "xCut/shop/service"
)
//
//func TestAdminIndex(t *testing.T) {
//	mux := http.NewServeMux()
//	admin := NewAdminDashboardHandler()
//	mux.HandleFunc("/handler/admin_dashboard_handler.go".)
//	ts := httptest.NewTLSServer(mux)
//	defer ts.Close()
//
//	tc := ts.Client()
//	URL := ts.URL
//	resp, err := tc.Login(URL+"/admin/", form)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
//	}
//
//}

func TestAdminBasicInfo(t *testing.T) {
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
	tmpl := template.Must(template.New("main").Funcs(fm).ParseGlob("../../../ui/templates/*"))

	shopRepo := shop.NewMockShopRepo(nil)
	shopService := shopServiceImport.NewShopService(shopRepo)

	adminHandler := NewAdminDashboardHandler(tmpl,shopService,nil,nil,nil,nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/basicInfo", adminHandler.AdminBasicInfo)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	URL := ts.URL

	resp, err := tc.Get(URL + "/admin/basicInfo")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("Mock Shop 01")) {
		t.Errorf("want body to contain %q", body)
	}

}

