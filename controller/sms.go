package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// GetSMSGatewayAndTemplate -  Return SMS Gateway & Template Details
func GetSMSGatewayAndTemplate(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - SMS - GetSMSGatewayAndTemplate")
	defer util.CommonDeferred(w, r, "Controller", "SMS	", "GetSMSGatewayAndTemplate")
	var retMap = make(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	Gateway, err := model.GetSmsGatewayDetail(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	Template, err := model.GetSmsTemplate(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	retMap["gateway"] = Gateway
	retMap["template"] = Template

	util.RespondData(r, w, retMap, 200)
}

// SMSTemplateListing - Return Datatable Listing Of SMS Template
func SMSTemplateListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - SMSTemplateListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "SMSTemplateListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.SMSTemplateListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetSMSTemplate -  Return Email Template Details
func GetSMSTemplate(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - SMS - GetSMSTemplate")
	defer util.CommonDeferred(w, r, "Controller", "SMS", "GetSMSTemplate")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetSmsTemplate(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetSmsGatewayDetail -  Return Email Template Details
func GetSmsGatewayDetail(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - SMS - GetSmsGatewayDetail")
	defer util.CommonDeferred(w, r, "Controller", "SMS", "GetSmsGatewayDetail")

	retMap, err := model.GetSmsGatewayDetail(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// UpdateSMSTemplate -  Update SMS Details
func UpdateSMSTemplate(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - SMS - UpdateSMSTemplate")
	defer util.CommonDeferred(w, r, "Controller", "SMS", "UpdateSMSTemplate")
	var reqMap data.SMSTemplate
	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.ID = id

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(id, reqMap.TemplateName, reqMap.TemplateVariable, reqMap.Template, reqMap.Subject)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.UpdateSMSTemplate(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, nil, 204)
}

// UpdateSMSGateway -  Update SMS Details
func UpdateSMSGateway(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - SMS - UpdateSMSGateway")
	defer util.CommonDeferred(w, r, "Controller", "SMS", "UpdateSMSGateway")
	var reqMap data.SmsGateway
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.EndPoint, reqMap.Username, reqMap.Password)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.UpdateSMSGateway(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, nil, 204)
}
