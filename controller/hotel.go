package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// AddHotelInfo - Add Hotel Info
func AddHotelInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - AddHotelInfo")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "AddHotelInfo")
	var reqMap data.Hotel

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	// 2020-06-18 - HK - PropertyType Added
	// 2020-06-15 - HK - Locality Added
	ValidateString := ValidateNotNullStructString(reqMap.Name, reqMap.Manager, reqMap.ShortAddress, reqMap.Username, reqMap.Password, reqMap.Email, reqMap.PhoneCode1, reqMap.Phone1, reqMap.AccountManager, reqMap.Locality, reqMap.PropertyType)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.Latitude, reqMap.Longitude, reqMap.HotelStar, reqMap.City, reqMap.State, reqMap.Country)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "HOTEL", nil, map[string]string{"hotel_name": reqMap.Name}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	cnt, err = model.CheckDuplicateRecords(r, "HOTEL_CLIENT", nil, map[string]string{"username": reqMap.Username}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10006")
		}
		return
	}

	flag := model.AddHotelInfo(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// HotelListing - Datatable Hotel listing with filter and order
func HotelListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - HotelListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "HotelListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.HotelListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// UpdateCommissionSetting - UpdateCommission Setting
func UpdateCommissionSetting(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - UpdateCommissionSetting")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "UpdateCommissionSetting")
	var reqMap data.HotelCommission

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.HotelID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.HotelID)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.CommissionAmount, reqMap.AutoSetDeal, reqMap.AutoSetDays)
	if reqMap.AutoSetDeal == 1 && reqMap.AutoSetDays == 0 {
		util.Respond(r, w, nil, 406, "100012")
		return
	}
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.UpdateCommissionSetting(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UpdateAccountManager - Update Assign Manager
func UpdateAccountManager(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - UpdateAccountManager")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "UpdateAccountManager")
	var reqMap data.Hotel

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.AccountManager, reqMap.ID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.UpdateAccountManager(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UpdateHotelInfo - Update Hotel Info
func UpdateHotelInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - UpdateHotelInfo")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "UpdateHotelInfo")
	var reqMap data.Hotel

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	// 2020-06-18 - HK - PropertyType Added
	// 2020-06-15 - HK - Locality Added
	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Name, reqMap.Manager, reqMap.ShortAddress, reqMap.Email, reqMap.PhoneCode1, reqMap.Phone1, reqMap.Locality, reqMap.PropertyType)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.Latitude, reqMap.Longitude, reqMap.HotelStar, reqMap.City, reqMap.State, reqMap.Country)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "HOTEL", nil, map[string]string{"hotel_name": reqMap.Name}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateHotelInfo(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// ResetHotelUserPwd - Reset Password Of Hotel User
func ResetHotelUserPwd(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - ResetHotelUserPwd")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "ResetHotelUserPwd")
	var reqMap data.Hotel

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.Password, reqMap.ID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.ResetHotelUserPwd(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetHotelInfo - Get Hotel Info Pass By ID
func GetHotelInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - GetHotelInfo")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "GetHotelInfo")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag, err := model.GetHotelInfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, flag, 200, "")
}

// ViewHotelInfo - View Hotel Info Pass By ID
func ViewHotelInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - ViewHotelInfo")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "ViewHotelInfo")
	var Panel = context.Get(r, "Side").(string)
	var id string
	//Handle Admin & Partner Panel API
	if Panel == "TP-BACKOFFICE" || Panel == "TP-PARTNER" {
		vars := mux.Vars(r)
		id = vars["id"]
	}
	// } else if Panel == "TP-PARTNER" {
	// 	id = context.Get(r, "HotelId").(string)
	// }

	//check of property associated with partner or not
	if !model.IsPartnerContainProperty(r, id, "") {
		util.Respond(r, w, nil, 406, "")
		return
	}

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag, err := model.ViewHotelInfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, flag, 200, "")
}

// UpdateBankDetails - Update Bank Details Setting
func UpdateBankDetails(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - UpdateBankDetails")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "UpdateBankDetails")
	var reqMap data.HotelBankDetails
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	var Panel = context.Get(r, "Side").(string)
	//Handle Admin & Partner Panel API
	if Panel == "TP-BACKOFFICE" || Panel == "TP-PARTNER" {
		vars := mux.Vars(r)
		reqMap.HotelID = vars["id"]
	}

	//check of property associated with partner or not
	if !model.IsPartnerContainProperty(r, reqMap.HotelID, "") {
		util.Respond(r, w, nil, 406, "")
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.HotelID, reqMap.AccountNumber, reqMap.SwiftCode, reqMap.AccountHolderName, reqMap.Bank)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.UpdateBankDetails(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetHotelListForOtherModule - Get Property Type List For Other Module
func GetHotelListForOtherModule(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - GetHotelListForOtherModule")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "GetHotelListForOtherModule")

	retMap, err := model.GetHotelListForOtherModule(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// UpdateHotelStatusToLive - Updates Hotel Status As Live And It Gets Listed In Mongo Search
func UpdateHotelStatusToLive(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - UpdateHotelStatusToLive")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "UpdateHotelStatusToLive")
	var reqMap data.StatusForSingle

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.ID = id
	var status = int(reqMap.Status)

	if reqMap.ID == "" {
		util.RespondBadRequest(r, w)
		return
	}

	if status != 1 && status != 2 {
		util.RespondBadRequest(r, w)
		return
	}

	// Check whether the Hotel Data is completely filled or not
	hotelVerificationFlg, errCodeString := model.CheckHotelVerificationEligibility(reqMap.ID)
	if !hotelVerificationFlg {
		util.Respond(r, w, nil, 406, errCodeString)
		return
	}

	flag, err := model.UpdateHotelStatusToLive(r, status, reqMap.ID)
	if !flag || err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, nil, 204)
}

// HotelierListing - Datatable HotelierListing listing with filter and order
func HotelierListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - HotelierListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "HotelierListing")
	var reqMap data.JQueryTableUI
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.HotelierListing(r, reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// ApprovedHotel - Approved hotel function
func ApprovedHotel(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - ApprovedHotel")
	defer util.CommonDeferred(w, r, "Controller", "V_Hotel", "ApprovedHotel")
	var reqMap data.ApprovedHotel
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.HotelID = id

	validateString := ValidateNotNullStructString(reqMap.HotelID)
	if validateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// Check whether the Hotel Data is completely filled or not
	hotelVerificationFlg, errCodeString := model.CheckHotelVerificationEligibility(reqMap.HotelID)
	if !hotelVerificationFlg {
		util.Respond(r, w, nil, 406, errCodeString)
		return
	}

	flag := model.ApprovedHotel(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UpdateBooking - Update booking status
func UpdateBooking(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Hotel - UpdateBooking")
	defer util.CommonDeferred(w, r, "Controller", "V_Booking", "UpdateBooking")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	var reqMap data.Status
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	errBooking := model.UpdateBooking(r, id, reqMap)
	if errBooking != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}
