package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddListingInquiry - Adds Property Listing Inquiry
func AddListingInquiry(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Inquiry - AddListingInquiry")
	defer util.CommonDeferred(w, r, "Controller", "Inquiry", "AddListingInquiry")

	var reqMap data.HotelInquiryModel

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.FirstName, reqMap.LastName, reqMap.Email, reqMap.PropertyName, reqMap.Address)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.PhoneCode, reqMap.Phone, reqMap.City, reqMap.State, reqMap.Country, reqMap.ZipCode)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "HOTEL_INQUIRY", map[string]string{"email": reqMap.Email, "property_name": reqMap.PropertyName}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddListingInquiry(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// OnBoardInquiryListing - Datatable On Board Inquiry Listing Data
func OnBoardInquiryListing(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Inquiry - OnBoardInquiryListing")
	defer util.CommonDeferred(w, r, "Controller", "Inquiry", "OnBoardInquiryListing")

	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.OnBoardInquiryListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetInquiryDetailInfo - Gets Hotel Inquiry Info
func GetInquiryDetailInfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Inquiry - GetInquiryDetailInfo")
	defer util.CommonDeferred(w, r, "Controller", "Inquiry", "GetInquiryDetailInfo")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetInquiryDetailInfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
