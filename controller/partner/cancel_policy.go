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

// AddCancelPolicy - Add Cancel Policy For hotel
func AddCancelPolicy(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Cancel_Policy - AddCancelPolicy")
	defer util.CommonDeferred(w, r, "Controller", "V_Cancel_Policy", "AddCancelPolicy")
	var reqMap data.CancelPolicy
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := controller.ValidateNotNullStructString(reqMap.Policy, reqMap.HotelID)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.IsNonRefundable, reqMap.BeforeDayCharge, reqMap.AfterDayCharge)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "CANCEL_POLICY", map[string]string{"policy_name": reqMap.Policy, "hotel_id": reqMap.HotelID}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := partner.AddCancelPolicy(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateCancelPolicy - Update Cancel Policy
func UpdateCancelPolicy(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Cancel_Policy - UpdateCancelPolicy")
	defer util.CommonDeferred(w, r, "Controller", "V_Cancel_Policy", "UpdateCancelPolicy")

	var reqMap data.CancelPolicy

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]
	// reqMap.HotelID = context.Get(r, "HotelId").(string)
	ValidateString := controller.ValidateNotNullStructString(reqMap.ID, reqMap.Policy, reqMap.HotelID)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.IsNonRefundable, reqMap.BeforeDayCharge, reqMap.AfterDayCharge)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "CANCEL_POLICY", map[string]string{"policy_name": reqMap.Policy, "hotel_id": reqMap.HotelID}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := partner.UpdateCancelPolicy(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// CancelPolicyListing - Return Datatable Listing Of Cancel Policy
func CancelPolicyListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Cancel_Policy - CancelPolicyListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Cancel_Policy", "CancelPolicyListing")
	var reqMap data.JQueryTableUI
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	//Validation for property
	if !model.IsPartnerContainProperty(r, reqMap.HotelID, "") {
		util.Respond(r, w, nil, 406, "")
		return
	}

	Data, err := partner.CancelPolicyListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetCancelPolicyList - Get Cancel Policy List For Other Module
func GetCancelPolicyList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Cancel_Policy - GetCancelPolicyList")
	defer util.CommonDeferred(w, r, "Controller", "V_Cancel_Policy", "GetCancelPolicyList")
	vars := mux.Vars(r)
	HotelID := vars["hotelid"]
	retMap, err := partner.GetCancelPolicyList(r, HotelID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetCancelPolicyInfo - Get Cancel Policy Info
func GetCancelPolicyInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Cancel_Policy - GetCancelPolicyInfo")
	defer util.CommonDeferred(w, r, "Controller", "V_Cancel_Policy", "GetCancelPolicyInfo")

	vars := mux.Vars(r)
	ID := vars["id"]
	HotelID := r.URL.Query().Get("hotelid")

	ValidateString := controller.ValidateNotNullStructString(ID, HotelID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := partner.GetCancelPolicyInfo(r, ID, HotelID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
