package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddLocality - Add Locality
func AddLocality(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Locality - AddLocality")
	defer util.CommonDeferred(w, r, "Controller", "V_Locality", "AddLocality")

	var reqMap data.Locality

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Locality)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.City)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "LOCALITY", map[string]string{"locality": reqMap.Locality, "city_id": fmt.Sprintf("%.0f", reqMap.City)}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddLocality(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateLocality - Add Locality
func UpdateLocality(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Locality - UpdateLocality")
	defer util.CommonDeferred(w, r, "Controller", "V_Locality", "UpdateLocality")

	var reqMap data.Locality

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Locality)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.City)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "LOCALITY", map[string]string{"locality": reqMap.Locality, "city_id": fmt.Sprintf("%.0f", reqMap.City)}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateLocality(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetLocality - Get Locality Detail By ID
func GetLocality(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Locality - GetLocality")
	defer util.CommonDeferred(w, r, "Controller", "V_Locality", "GetLocality")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetLocality(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// LocalityListing - Return Datatable Listing Of Locality
func LocalityListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Locality - LocalityListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Locality", "LocalityListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.LocalityListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
