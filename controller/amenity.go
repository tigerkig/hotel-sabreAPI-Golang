package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddAmenity - Add Amenity
func AddAmenity(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - AddAmenity")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "AddAmenity")
	var reqMap data.Amenity

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Type, reqMap.Name, reqMap.Icon)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.IsStarAmenity, reqMap.AmenityOf)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// cnt, err := model.CheckDuplicateRecords(r, "AMENITY", nil, map[string]string{"name": reqMap.Name}, "0")
	cnt, err := model.CheckDuplicateRecords(r, "AMENITY", map[string]string{"name": reqMap.Name, "amenity_type_id": reqMap.Type}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddAmenity(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateAmenity - Update Amenity
func UpdateAmenity(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - UpdateAmenity")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "UpdateAmenity")
	var reqMap data.Amenity

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Type, reqMap.Name, reqMap.Icon)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.IsStarAmenity, reqMap.AmenityOf)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// cnt, err := model.CheckDuplicateRecords(r, "AMENITY", nil, map[string]string{"name": reqMap.Name}, reqMap.ID)
	cnt, err := model.CheckDuplicateRecords(r, "AMENITY", map[string]string{"name": reqMap.Name, "amenity_type_id": reqMap.Type}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateAmenity(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetAmenity - Get Amenity Detail By ID
func GetAmenity(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - GetAmenity")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "GetAmenity")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetAmenity(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetAmenityList -  Get Amenity Active List
func GetAmenityList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - GetAmenityList")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "GetAmenityList")

	retMap, err := model.GetAmenityList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// AmenityListing - Return Datatable Listing Of Amenity
func AmenityListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - AmenityListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "AmenityListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.AmenityListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// AmenityTypeWiseAmenity - Return amenity type wise amenity data
func AmenityTypeWiseAmenity(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - AmenityTypeWiseAmenity")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "AmenityTypeWiseAmenity")
	//var HotelID = context.Get(r, "HotelId").(string)
	HotelID := r.URL.Query().Get("hotelid")
	if HotelID == "" {
		util.RespondBadRequest(r, w)
		return
	}
	retMap, err := model.AmenityTypeWiseAmenity(r, HotelID)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetAmenityListV1 -  Get Amenity Active List
func GetAmenityListV1(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - GetAmenityListV1")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "GetAmenityListV1")

	vars := mux.Vars(r)
	CatgID := vars["id"]

	if CatgID == "" || CatgID == "0" {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetAmenityListV1(r, CatgID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// AmenityTypeWiseAmenityForRoom - Return amenity type wise amenity data for room
func AmenityTypeWiseAmenityForRoom(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Amenity - AmenityTypeWiseAmenityForRoom")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "AmenityTypeWiseAmenityForRoom")

	// var HotelID = context.Get(r, "HotelId").(string)

	vars := mux.Vars(r)
	roomID := vars["id"]

	HotelID := r.URL.Query().Get("hotelid")
	ValidateString := ValidateNotNullStructString(roomID, HotelID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.AmenityTypeWiseAmenityForRoom(r, HotelID, roomID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
