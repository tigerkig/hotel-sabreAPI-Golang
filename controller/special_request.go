package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddSpecialRequest - Add Special Request
func AddSpecialRequest(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Special_Request - AddSpecialRequest")
	defer util.CommonDeferred(w, r, "Controller", "V_Special_Request", "AddSpecialRequest")
	var reqMap data.SpecialRequest

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Request)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "SPECIAL_REQUEST", nil, map[string]string{"special_request": reqMap.Request}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddSpecialRequest(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateSpecialRequest - Update Sepcial Request
func UpdateSpecialRequest(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Special_Request - UpdateSpecialRequest")
	defer util.CommonDeferred(w, r, "Controller", "V_Special_Request", "UpdateSpecialRequest")
	var reqMap data.SpecialRequest

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Request)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "SPECIAL_REQUEST", nil, map[string]string{"special_request": reqMap.Request}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateSpecialRequest(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// SpecialRequestListing - Datatable Special Request listing with filter and order
func SpecialRequestListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Special_Request - SpecialRequestListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Special_Request", "SpecialRequestListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.SpecialRequestListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetSpecialRequest - This function returns special request which are active. Function made for web.
func GetSpecialRequest(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Special_Request - GetSpecialRequest")
	defer util.CommonDeferred(w, r, "Controller", "Special_Request", "GetSpecialRequest")

	Data, flg := model.GetSpecialRequest(r)
	if !flg {
		util.RespondWithError(r, w, "500")
		return
	}
	util.Respond(r, w, Data, 200, "")
}
