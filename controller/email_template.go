package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// UpdateEmailTemplate - Update Email Template
func UpdateEmailTemplate(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Email_Template - UpdateEmailConfig")
	defer util.CommonDeferred(w, r, "Controller", "Email_Template", "UpdateEmailConfig")
	var reqMap data.EmailTemplate

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Email, reqMap.Template, reqMap.Subject)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.UpdateEmailTemplate(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetEmailTemplate -  Return Email Template Details
func GetEmailTemplate(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Email_Template - GetEmailTemplate")
	defer util.CommonDeferred(w, r, "Controller", "Email_Template", "GetEmailTemplate")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetEmailTemplate(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// EmailTemplateListing - Return Datatable Listing Of Email Template
func EmailTemplateListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Amenity - EmailTemplateListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Amenity", "EmailTemplateListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.EmailTemplateListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetEmailTemplateDetailInfo -  Return Email Template Details
func GetEmailTemplateDetailInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Email_Template - GetEmailTemplateDetailInfo")
	defer util.CommonDeferred(w, r, "Controller", "Email_Template", "GetEmailTemplateDetailInfo")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetEmailTemplateDetailInfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetEmailList -  Return Email List
func GetEmailList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Email_Template - GetEmailList")
	defer util.CommonDeferred(w, r, "Controller", "Email_Template", "GetEmailList")
	retMap, err := model.GetEmailList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	resMap := make(map[string]interface{})
	if len(retMap) > 0 {
		resMap["data"] = retMap
	} else {
		resMap["data"] = []string{}
	}

	util.RespondData(r, w, resMap, 200)
}
