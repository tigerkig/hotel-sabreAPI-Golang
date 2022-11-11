package controller

import (
	"encoding/json"
	"net/http"
	"tp-system/model"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/mux"
)

// AddPropPertyTag - Add Property Tag
func AddPropPertyTag(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Property_Tag - AddPropPertyTag")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Tag", "AddPropPertyTag")
	var reqMap data.PropertyTags

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Tag)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "PROPERTY_TAG", nil, map[string]string{"tag": reqMap.Tag}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddPropertyTag(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdatePropPertyTag - Update Property Tag
func UpdatePropPertyTag(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Property_Tag - UpdatePropPertyTag")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Tag", "UpdatePropPertyTag")
	var reqMap data.PropertyTags

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Tag)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "PROPERTY_TAG", nil, map[string]string{"tag": reqMap.Tag}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdatePropPertyTag(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetPropPertyTagList - Get Property Tag List For Other Module
func GetPropPertyTagList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Property_Tag - GetPropPertyTagList")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Tag", "GetPropPertyTagList")
	retMap, err := model.GetPropertyTagList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// PropPertyTagListing - Datatable Property Tag listing with filter and order
func PropPertyTagListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Property_Tag - PropPertyTagListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Tag", "PropPertyTagListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.PropPertyTagListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
