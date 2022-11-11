package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddExtraBedType - Add Extra Bed Type
func AddExtraBedType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Extra_Bed_Type - AddExtraBedType")
	defer util.CommonDeferred(w, r, "Controller", "V_Extra_Bed_Type", "AddExtraBedType")
	var reqMap data.ExtraBedType

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Type)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "EXTRA_BED_TYPE", nil, map[string]string{"extra_bed_name": reqMap.Type}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddExtraBedType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateExtraBedType - Update Extra Bed Type
func UpdateExtraBedType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Extra_Bed_Type - UpdateExtraBedType")
	defer util.CommonDeferred(w, r, "Controller", "V_Extra_Bed_Type", "UpdateExtraBedType")
	var reqMap data.ExtraBedType

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Type)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "EXTRA_BED_TYPE", nil, map[string]string{"extra_bed_name": reqMap.Type}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateExtraBedType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetExtraBedTypeList - Get Extra Bed Type List For Other Module
func GetExtraBedTypeList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Extra_Bed_Type - GetExtraBedTypeList")
	defer util.CommonDeferred(w, r, "Controller", "V_Extra_Bed_Type", "GetExtraBedTypeList")
	retMap, err := model.GetExtraBedTypeList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// ExtraBedTypeListing - Datatable Extra Bed Type listing with filter and order
func ExtraBedTypeListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Extra_Bed_Type - ExtraBedTypeListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Extra_Bed_Type", "ExtraBedTypeListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.ExtraBedTypeListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
