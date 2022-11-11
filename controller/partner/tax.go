package partner

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/controller"
	"tp-system/model"
	"tp-system/model/partner"

	"github.com/gorilla/mux"
)

// AddTax - Add Tax For Hotel
func AddTax(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Tax - AddTax")
	defer util.CommonDeferred(w, r, "Controller", "V_Tax", "AddTax")

	var reqMap data.Tax

	// Converts JSON into a Go value.
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	// Validates Incoming Request Data According To Tax Model Defined In tp-api-common Tax Struct
	ValidateString := controller.ValidateNotNullStructString(reqMap.Name, reqMap.Description, reqMap.Type, reqMap.HotelID)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.Amount)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// IF TaxType IS PERCENTAGE Then Value Must Not Be >100 And <0
	if reqMap.Type == "PERCENTAGE" {
		if reqMap.Amount > 100 || reqMap.Amount <= 0 {
			util.RespondBadRequest(r, w)
			return
		}
	}

	HotelID := reqMap.HotelID
	// checks if hotel id exists in system with optional status check
	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, HotelID, true, false)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}

	// // Retrieves Hotel ID From Context
	// HotelID := context.Get(r, "HotelId")
	// // To Check Whether Same Tax Name Exists With The Same Hotel ID Or Not
	// cnt, err := model.CheckDuplicateRecords(r, "TAX", map[string]string{"tax": reqMap.Name, "hotel_id": HotelID.(string)}, nil, "0")
	cnt, err := model.CheckDuplicateRecords(r, "TAX", map[string]string{"tax": reqMap.Name, "hotel_id": HotelID}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	// Calls Partner Model Function To Insert Tax Details
	flag := partner.AddTax(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateTax - Update Tax For Hotel
func UpdateTax(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Tax - UpdateTax")
	defer util.CommonDeferred(w, r, "Controller", "V_Tax", "UpdateTax")

	var reqMap data.Tax

	// Converts JSON into a Go value.
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"] // Gets ID FROM URL

	// Validates Incoming Request Data According To Tax Model Defined In tp-api-common Tax Struct
	// ValidateString := controller.ValidateNotNullStructString(reqMap.Name, reqMap.Description, reqMap.Type, reqMap.ID)
	ValidateString := controller.ValidateNotNullStructString(reqMap.Name, reqMap.Description, reqMap.Type, reqMap.ID, reqMap.HotelID)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.Amount)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// IF TaxType IS PERCENTAGE Then Value Must Not Be >100 And <0
	if reqMap.Type == "PERCENTAGE" {
		if reqMap.Amount > 100 || reqMap.Amount <= 0 {
			util.RespondBadRequest(r, w)
			return
		}
	}

	// // Retrieves Hotel ID From Context
	// HotelID := context.Get(r, "HotelId")
	// // To Check Whether Same Tax Name Exists With The Same Hotel ID Or Not By Passing ID Here As Not To Compare With Own ID
	// cnt, err := model.CheckDuplicateRecords(r, "TAX", map[string]string{"tax": reqMap.Name, "hotel_id": HotelID.(string)}, nil, reqMap.ID)
	HotelID := reqMap.HotelID
	// checks if hotel id exists in system with optional status check
	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, HotelID, true, false)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "TAX", map[string]string{"tax": reqMap.Name, "hotel_id": HotelID}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	// Calls Partner Model Function To Update Tax Details
	flag := partner.UpdateTax(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetTaxInfo - Get Tax Info
func GetTaxInfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Tax - GetTaxInfo")
	defer util.CommonDeferred(w, r, "Controller", "V_Tax", "GetTaxInfo")

	vars := mux.Vars(r)
	ID := vars["id"] // Gets ID FROM URL

	HotelID := r.URL.Query().Get("hotelid")

	//Validates If ID Passed Is Blank Or NOT
	ValidateString := controller.ValidateNotNullStructString(ID, HotelID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// Calls Partner Model Function To Get Tax Details
	retMap, err := partner.GetTaxInfo(r, ID, HotelID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// TaxListing - Return Datatable Listing Of Tax
func TaxListing(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Tax - TaxListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Tax", "TaxListing")

	var reqMap data.JQueryTableUI

	// Converts JSON into a Go value.
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	if reqMap.HotelID == "" {
		util.LogIt(r, "Controller - V_Tax - TaxListing - Hotel Id Missing")
		util.RespondBadRequest(r, w)
		return
	}

	// Calls Partner Model Function To Get Tax Listing
	Data, err := partner.TaxListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
