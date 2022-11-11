package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// UpdatePaymentConfiguration - Update Payment Configuration
func UpdatePaymentConfiguration(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Payment_Configuration - UpdatePaymentConfiguration")
	defer util.CommonDeferred(w, r, "Controller", "Payment_Configuration", "UpdatePaymentConfiguration")
	var reqMap data.PaymentGatewayConfiguration

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.PaymentType, reqMap.URL, reqMap.AuthKey, reqMap.SecretKey, reqMap.Country)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.UpdatePaymentConfiguration(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetPaymentConfigDetail -  Return Payment Config Details
func GetPaymentConfigDetail(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Payment_Configuration - GetPaymentConfigDetail")
	defer util.CommonDeferred(w, r, "Controller", "Payment_Configuration", "GetPaymentConfigDetail")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetPaymentConfigDetail(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// ActivatePaymentGateway -  Activate Payment Gateway
func ActivatePaymentGateway(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Payment_Configuration - ActivatePaymentGateway")
	defer util.CommonDeferred(w, r, "Controller", "Payment_Configuration", "ActivatePaymentGateway")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 || id == "undefined" {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.ActivatePaymentGateway(r, id)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, nil, 204)
}

// GetActivatePaymentGateway -  Return Activate Payment Gateway
func GetActivatePaymentGateway(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Payment_Configuration - GetActivatePaymentGateway")
	defer util.CommonDeferred(w, r, "Controller", "Payment_Configuration", "GetActivatePaymentGateway")
	// Get active payment gateway data
	// Only one payment gateway at a time can we get
	Data, err := model.GetActivatePaymentGateway(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, Data, 200)
}

// PaymentGatewayListing - Return Datatable Listing Of Payment Gateway Configuration
func PaymentGatewayListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Payment_Configuration - PaymentGatewayListing")
	defer util.CommonDeferred(w, r, "Controller", "Payment_Configuration", "PaymentGatewayListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.PaymentGatewayListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetPaymentGatewayList -  Return Payment Config Details List
func GetPaymentGatewayList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Payment_Configuration - GetPaymentGatewayList")
	defer util.CommonDeferred(w, r, "Controller", "Payment_Configuration", "GetPaymentGatewayList")
	retMap, err := model.GetPaymentGatewayList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
