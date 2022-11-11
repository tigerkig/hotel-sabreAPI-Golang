package controller

import (
	"fmt"
	"net/http"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// CompanyInfo - Company Info List
func CompanyInfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Company - CompanyInfo")
	defer util.CommonDeferred(w, r, "Controller", "Company", "CompanyInfo")

	Data, err := model.GetCompanyInfo(r, "")
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// UpdateCompanyInfo - Updates Company Info
func UpdateCompanyInfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Company - UpdateCompanyInfo")
	defer util.CommonDeferred(w, r, "Controller", "Company", "UpdateCompanyInfo")

	var reqMap = make(map[string]interface{})

	vars := mux.Vars(r)
	ID := vars["id"]
	Name := r.FormValue("company_name")
	County := r.FormValue("country_id")
	State := r.FormValue("state_id")
	City := r.FormValue("city_id")
	ZipCode := r.FormValue("zip_code")
	Address := r.FormValue("address")
	RegAddress := r.FormValue("registered_office_address")

	if ID == "" || Name == "" || County == "" || State == "" || City == "" || ZipCode == "" || Address == "" || RegAddress == "" {
		util.RespondBadRequest(r, w)
		return
	}

	reqMap["company_name"] = Name
	reqMap["country_id"] = County
	reqMap["state_id"] = State
	reqMap["city_id"] = City
	reqMap["zip_code"] = ZipCode
	reqMap["address"] = Address
	reqMap["registered_office_address"] = RegAddress
	reqMap["id"] = ID

	imgfile, _, _ := r.FormFile("image")
	if imgfile != nil {
		ImageName, err := UploadSingleImageFormData(r, "company_logo", "image")
		if err != nil {
			util.LogIt(r, fmt.Sprint("Controller - Company - UpdateCompanyInfo - Error While Uploading Company Logo"))
			util.RespondBadRequest(r, w)
			return
		}
		reqMap["image"] = ImageName
	} else {
		reqMap["image"] = r.FormValue("image")
	}

	flag := model.UpdateCompanyInfo(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetCompanyInfo - Get Company Info By ID
func GetCompanyInfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Company - GetCompanyInfo")
	defer util.CommonDeferred(w, r, "Controller", "Company", "GetCompanyInfo")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetCompanyInfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
