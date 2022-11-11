package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddBedType - Add Bed Type
func AddBedType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Bed_Type - AddBedType")
	defer util.CommonDeferred(w, r, "Controller", "V_Bed_Type", "AddBedType")
	var reqMap data.BedTypes

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

	cnt, err := model.CheckDuplicateRecords(r, "BED_TYPE", nil, map[string]string{"bed_type": reqMap.Type}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddBedType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateBedType - Update Bed Type
func UpdateBedType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Bed_Type - UpdateBedType")
	defer util.CommonDeferred(w, r, "Controller", "V_Bed_Type", "UpdateBedType")
	var reqMap data.BedTypes

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

	cnt, err := model.CheckDuplicateRecords(r, "BED_TYPE", nil, map[string]string{"bed_type": reqMap.Type}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateBedType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetBedTypeList - Get Bed Type List For Other Module
func GetBedTypeList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Bed_Type - GetBedTypeList")
	defer util.CommonDeferred(w, r, "Controller", "V_Bed_Type", "GetBedTypeList")
	retMap, err := model.GetBedTypeList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// BedTypeListing - Datatable Bed Type listing with filter and order
func BedTypeListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Bed_Type - BedTypeListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Bed_Type", "BedTypeListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.BedTypeListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetBedType - Get Bed Type Detail By ID - 2021-04-21 - HK
func GetBedType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Bed_Type - GetBedType")
	defer util.CommonDeferred(w, r, "Controller", "V_Bed_Type", "GetBedType")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetBedType(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
