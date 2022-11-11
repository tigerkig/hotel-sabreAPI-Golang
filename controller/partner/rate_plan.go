package partner

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/controller"
	"tp-system/model"
	"tp-system/model/partner"

	"github.com/gorilla/mux"
)

func AddInv(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_RatePlan - AddInv")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_RatePlan", "AddInv")

	vars := mux.Vars(r)
	ID := vars["id"] // Gets ID FROM URL

	flag := partner.AddInv(r, ID)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 200, "")
}

// AddRatePlan - Add Rate Plan
func AddRatePlan(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_RatePlan - AddRatePlan")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_RatePlan", "AddRatePlan")

	var reqMap data.RatePlan

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	// ValidateString := controller.ValidateNotNullStructString(reqMap.Name, reqMap.RoomType, reqMap.CancelPolicy, reqMap.Inclusion, reqMap.MealPlan)
	ValidateString := controller.ValidateNotNullStructString(reqMap.Name, reqMap.RoomType, reqMap.CancelPolicy, reqMap.Inclusion, reqMap.MealPlan, reqMap.HotelID)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.IsPayAtHotel, reqMap.Rate, reqMap.SortOrder)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	HotelID := reqMap.HotelID
	// checks if hotel id exists in system with optional status check
	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, HotelID, false, false)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}

	// HotelID := context.Get(r, "HotelId")
	// cnt, err := model.CheckDuplicateRecords(r, "RATE_PLAN", map[string]string{"rate_plan_name": reqMap.Name, "room_type_id": reqMap.RoomType, "hotel_id": HotelID.(string)}, nil, "0")
	cnt, err := model.CheckDuplicateRecords(r, "RATE_PLAN", map[string]string{"rate_plan_name": reqMap.Name, "room_type_id": reqMap.RoomType, "hotel_id": HotelID}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := partner.AddRatePlan(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateRatePlan - Update Rate Plan
func UpdateRatePlan(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_RatePlan - UpdateRatePlan")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_RatePlan", "UpdateRatePlan")

	var reqMap data.RatePlan

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	// ValidateString := controller.ValidateNotNullStructString(reqMap.ID, reqMap.Name, reqMap.RoomType, reqMap.CancelPolicy, reqMap.Inclusion, reqMap.MealPlan)
	ValidateString := controller.ValidateNotNullStructString(reqMap.ID, reqMap.Name, reqMap.RoomType, reqMap.CancelPolicy, reqMap.Inclusion, reqMap.MealPlan, reqMap.HotelID)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.IsPayAtHotel, reqMap.Rate, reqMap.SortOrder)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// HotelID := context.Get(r, "HotelId")
	// cnt, err := model.CheckDuplicateRecords(r, "RATE_PLAN", map[string]string{"rate_plan_name": reqMap.Name, "room_type_id": reqMap.RoomType, "hotel_id": HotelID.(string)}, nil, reqMap.ID)

	HotelID := reqMap.HotelID
	// checks if hotel id exists in system with optional status check
	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, HotelID, true, false)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}
	cnt, err := model.CheckDuplicateRecords(r, "RATE_PLAN", map[string]string{"rate_plan_name": reqMap.Name, "room_type_id": reqMap.RoomType, "hotel_id": HotelID}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := partner.UpdateRatePlan(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetRatePlan - Get Rate Plan Info
func GetRatePlan(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_RatePlan - GetRatePlan")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_RatePlan", "GetRatePlan")

	vars := mux.Vars(r)
	ID := vars["id"]

	HotelID := r.URL.Query().Get("hotelid")

	ValidateString := controller.ValidateNotNullStructString(ID, HotelID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, HotelID, true, false)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}

	retMap, err := partner.GetRatePlan(r, ID, HotelID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// RatePlanListing - Return Datatable Listing Of Rate Plan
func RatePlanListing(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_RatePlan - RatePlanListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_RatePlan", "RatePlanListing")

	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	if reqMap.HotelID == "" {
		util.LogIt(r, "Controller - V_Partner_RatePlan - RatePlanListing - Hotel Id Missing")
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := partner.RatePlanListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetRatePlanList - Get Rate Plan List For Other Module
func GetRatePlanList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_RatePlan - GetRatePlanList")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_RatePlan", "GetRatePlanList")

	vars := mux.Vars(r)
	RoomID := vars["id"]
	HotelID := r.URL.Query().Get("hotelid")
	if RoomID == "" || RoomID == "0" || HotelID == "" {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := partner.GetRatePlanList(r, RoomID, HotelID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
