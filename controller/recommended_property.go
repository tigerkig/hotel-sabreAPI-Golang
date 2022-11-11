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

// AddHotelToRecommendedList - Adds Hotel Into Recommended List
func AddHotelToRecommendedList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Recommended_Hotel - AddHotelToRecommendedList")
	defer util.CommonDeferred(w, r, "Controller", "Recommended_Hotel", "AddHotelToRecommendedList")

	var reqMap data.RecommendedHotel

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.HotelID)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.SortOrder)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	TotalRecordCount, err := model.GetRecommendedPropertyCount(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	if TotalRecordCount >= 6 {
		util.Respond(r, w, nil, 406, "100014")
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "RECOMMENDED_HOTEL", nil, map[string]string{"hotel_id": reqMap.HotelID, "sort_order": fmt.Sprintf("%.0f", reqMap.SortOrder)}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddHotelToRecommendedList(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateHotelToRecommendedList - Updates Hotel Into Recommended List
func UpdateHotelToRecommendedList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Recommended_Hotel - UpdateHotelToRecommendedList")
	defer util.CommonDeferred(w, r, "Controller", "Recommended_Hotel", "UpdateHotelToRecommendedList")

	var reqMap data.RecommendedHotel

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.HotelID, reqMap.ID)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.SortOrder)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "RECOMMENDED_HOTEL", nil, map[string]string{"hotel_id": reqMap.HotelID, "sort_order": fmt.Sprintf("%.0f", reqMap.SortOrder)}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateHotelToRecommendedList(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// RecommendedHotelList - Datatable Listing Of Recommended Hotels
func RecommendedHotelList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Recommended_Hotel - RecommendedHotelList")
	defer util.CommonDeferred(w, r, "Controller", "Recommended_Hotel", "RecommendedHotelList")

	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.RecommendedHotelList(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetRecommendedHotelInfo - Gets Recommended Hotel Info
func GetRecommendedHotelInfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Recommended_Hotel - GetRecommendedHotelInfo")
	defer util.CommonDeferred(w, r, "Controller", "Recommended_Hotel", "GetRecommendedHotelInfo")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetRecommendedHotelInfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
