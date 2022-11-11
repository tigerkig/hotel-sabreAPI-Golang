package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddInclusion - Add Inclusion That's Free BreakFast etc
func AddInclusion(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Inclusion - AddInclusion")
	defer util.CommonDeferred(w, r, "Controller", "V_Inclusion", "AddInclusion")
	var reqMap data.Inclusion

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Inclusion)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "INCLUSION", nil, map[string]string{"inclusion": reqMap.Inclusion}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddInclusion(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateInclusion - Update Inclusion That's Free BreakFast etc
func UpdateInclusion(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Inclusion - UpdateInclusion")
	defer util.CommonDeferred(w, r, "Controller", "V_Inclusion", "UpdateInclusion")
	var reqMap data.Inclusion

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Inclusion)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "INCLUSION", nil, map[string]string{"inclusion": reqMap.Inclusion}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateInclusion(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetInclusionList - Get Inclusion List For Other Module
func GetInclusionList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Inclusion - GetInclusionList")
	defer util.CommonDeferred(w, r, "Controller", "V_Inclusion", "GetInclusionList")
	retMap, err := model.GetInclusionList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// InclusionListing - Datatable Inclusion listing with filter and order
func InclusionListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Inclusion - InclusionListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Inclusion", "InclusionListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.InclusionListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetInclusion - Get Inclusion Detail By ID - 2021-04-20 - HK
func GetInclusion(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Inclusion - GetInclusion")
	defer util.CommonDeferred(w, r, "Controller", "V_Inclusion", "GetInclusion")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetInclusion(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
