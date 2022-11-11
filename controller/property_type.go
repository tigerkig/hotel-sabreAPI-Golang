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

// AddPropertyType - Add Property Type
func AddPropertyType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Property_Type - AddPropertyType")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Type", "AddPropertyType")
	var reqMap data.PropertyType

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

	cnt, err := model.CheckDuplicateRecords(r, "PROPERTY_TYPE", nil, map[string]string{"type": reqMap.Type}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddPropertyType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// AddPropertyTypeNew - Add Property Type
func AddPropertyTypeNew(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Property_Type - AddPropertyTypeNew")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Type", "AddPropertyTypeNew")

	var reqMap = make(map[string]interface{})

	PropertyType := r.FormValue("type")
	reqMap["type"] = PropertyType

	cnt, err := model.CheckDuplicateRecords(r, "PROPERTY_TYPE", nil, map[string]string{"type": PropertyType}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	imgfile, _, err := r.FormFile("image")
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}
	if imgfile != nil {
		ImageName, err := UploadSingleImageFormData(r, "property_type", "image")
		if err != nil {
			util.RespondBadRequest(r, w)
			return
		}
		reqMap["image"] = ImageName
		defer imgfile.Close()
	}

	flag := model.AddPropertyTypeNew(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdatePropertyType - Update Property Type
func UpdatePropertyType(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Property_Type - UpdatePropertyType")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Type", "UpdatePropertyType")
	var reqMap data.PropertyType

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

	cnt, err := model.CheckDuplicateRecords(r, "PROPERTY_TYPE", nil, map[string]string{"type": reqMap.Type}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdatePropertyType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UpdatePropertyTypeNew - Update Property Type
func UpdatePropertyTypeNew(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Property_Type - UpdatePropertyTypeNew")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Type", "UpdatePropertyTypeNew")

	var reqMap = make(map[string]interface{})

	vars := mux.Vars(r)
	ID := vars["id"]

	PropertyType := r.FormValue("type")

	if ID == "" || PropertyType == "" {
		util.RespondBadRequest(r, w)
		return
	}

	reqMap["type"] = PropertyType
	reqMap["id"] = ID

	cnt, err := model.CheckDuplicateRecords(r, "PROPERTY_TYPE", nil, map[string]string{"type": PropertyType}, ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	imgfile, _, err := r.FormFile("image")
	if imgfile != nil {
		ImageName, err := UploadSingleImageFormData(r, "property_type", "image")
		if err != nil {
			util.LogIt(r, fmt.Sprint("Controller - V_Property_Type - UpdatePropertyTypeNew - Error While Updating Property Type Pic"))
			util.RespondBadRequest(r, w)
			return
		}
		reqMap["image"] = ImageName
	} else {
		reqMap["image"] = r.FormValue("image")
	}

	flag := model.UpdatePropertyTypeNew(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// PropertyTypeListing - Datatable Property type listing with filter and order
func PropertyTypeListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Property_Type - PropertyTypeListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Type", "PropertyTypeListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.PropertyTypeListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetPropertyTypeInfo - Get Property Type Info
func GetPropertyTypeInfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Property_Type - GetPropertyTypeInfo")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Type", "GetPropertyTypeInfo")

	vars := mux.Vars(r)
	ID := vars["id"]

	ValidateString := ValidateNotNullStructString(ID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetPropertyTypeInfo(r, ID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetPropertyTypeList - Get Property Type List For Other Module
func GetPropertyTypeList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Property_Type - GetPropertyTypeList")
	defer util.CommonDeferred(w, r, "Controller", "V_Property_Type", "GetPropertyTypeList")

	retMap, err := model.GetPropertyTypeList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
