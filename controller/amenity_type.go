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

// AddAmenityType - Add Amenity Type
func AddAmenityType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity_Type - AddAmenityType")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity_Type", "AddAmenityType")
	var reqMap data.AmenityType

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Type)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.AmenityOf)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// cnt, err := model.CheckDuplicateRecords(r, "AMENITY_TYPE", nil, map[string]string{"type": reqMap.Type}, "0")
	cnt, err := model.CheckDuplicateRecords(r, "AMENITY_TYPE", map[string]string{"type": reqMap.Type, "amenity_of": fmt.Sprintf("%.0f", reqMap.AmenityOf)}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddAmenityType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateAmenityType - Update Amenity Type
func UpdateAmenityType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity_Type - UpdateAmenityType")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity_Type", "UpdateAmenityType")
	var reqMap data.AmenityType

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Type)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.AmenityOf)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// cnt, err := model.CheckDuplicateRecords(r, "AMENITY_TYPE", nil, map[string]string{"type": reqMap.Type}, reqMap.ID)
	cnt, err := model.CheckDuplicateRecords(r, "AMENITY_TYPE", map[string]string{"type": reqMap.Type, "amenity_of": fmt.Sprintf("%.0f", reqMap.AmenityOf)}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateAmenityType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetAmenityTypeList - Get Amenity Type List For Other Module
func GetAmenityTypeList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity_Type - GetAmenityTypeList")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity_Type", "GetAmenityTypeList")
	retMap, err := model.GetAmenityTypeList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// AmenityTypeListing - Datatable amenity type listing with filter and order
func AmenityTypeListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity_Type - AmenityTypeListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity_Type", "AmenityTypeListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.AmenityTypeListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetAmenityTypeListCatgWise - Get Amenity Type List For Other Module
func GetAmenityTypeListCatgWise(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Amenity_Type - GetAmenityTypeListCatgWise")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity_Type", "GetAmenityTypeListCatgWise")

	vars := mux.Vars(r)
	ID := vars["id"]

	ValidateString := ValidateNotNullStructString(ID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetAmenityTypeListCatgWise(r, ID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetAmenityType - Get Amenity Type Info
func GetAmenityType(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Amenity_Type - GetAmenityType")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity_Type", "GetAmenityType")

	vars := mux.Vars(r)
	ID := vars["id"]

	ValidateString := ValidateNotNullStructString(ID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetAmenityType(r, ID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetAmenityTypeListV1 - Get Amenity Type List For Other Module
func GetAmenityTypeListV1(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Amenity_Type - GetAmenityTypeListV1")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity_Type", "GetAmenityTypeListV1")

	vars := mux.Vars(r)
	CatgID := vars["id"]

	if CatgID == "" || CatgID == "0" {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetAmenityTypeListV1(r, CatgID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
