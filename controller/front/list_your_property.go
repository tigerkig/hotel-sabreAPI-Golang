package front

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/controller"
	"tp-system/model"
	"tp-system/model/front"

	"github.com/gorilla/mux"
)

// AddListProperty - Adds Property Listing Inquiry
func AddListProperty(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Front-list-your-property - AddListProperty")
	defer util.CommonDeferred(w, r, "Controller", "Front-list-your-property", "AddListProperty")
	var reqMap data.ListYourProperty
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := controller.ValidateNotNullStructString(reqMap.FirstName, reqMap.LastName, reqMap.Email, reqMap.Password)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.PhoneCode, reqMap.Phone)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "LIST_YOUR_PROPERTY", map[string]string{"email": reqMap.Email}, nil, "0")
	if err != nil || cnt != 0 {
		if util.CheckErrorLog(r, err) {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	cnt, err = model.CheckDuplicateRecords(r, "HOTEL_CLIENT", map[string]string{"email": reqMap.Email, "username": reqMap.Email}, nil, "0")
	if err != nil || cnt != 0 {
		if util.CheckErrorLog(r, err) {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	if !front.AddListingInquiry(r, reqMap) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// VerifyPartnerActivateToken - Verify partner token
func VerifyPartnerActivateToken(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Front-list-your-property - VerifyPartnerActivateToken")
	defer util.CommonDeferred(w, r, "Controller", "Front-list-your-property", "VerifyPartnerActivateToken")
	var reqMap data.VerifyPartnerToken
	vars := mux.Vars(r)
	reqMap.Token = vars["token"]

	switch front.CheckPartnerAlreadyRegister(r, reqMap.Token, "TOKEN") {
	case 1, 4:
		util.LogIt(r, fmt.Sprintf("User link not expired"))
	case 2:
		util.Respond(r, w, nil, 409, "")
		return
	case 3, 5:
		util.Respond(r, w, nil, 406, "Link expired")
		return
	}

	if !front.VerifyJWTToken(r, reqMap) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 200, "")
}

// ListYourPropertyListing - Datatable ListYourProperty listing with filter and order
func ListYourPropertyListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Front-list-your-property - ListYourPropertyListing")
	defer util.CommonDeferred(w, r, "Controller", "Front-list-your-property", "ListYourPropertyListing")
	var reqMap data.JQueryTableUI
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := front.ListYourPropertyListing(r, reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
