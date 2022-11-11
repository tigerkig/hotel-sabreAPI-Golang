package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddEmailConfig - Adds Email Configuration
func AddEmailConfig(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Email - AddEmailConfig")
	defer util.CommonDeferred(w, r, "Controller", "Email", "AddEmailConfig")
	var reqMap data.EmailConfig

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.EmailID, reqMap.EmailName, reqMap.SMTPHost, reqMap.SMTPUser, reqMap.SMTPPassword, reqMap.Signature)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.SMTPPort)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "EMAIL_CONFIG", nil, map[string]string{"email_id": reqMap.EmailID}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddEmailConfig(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateEmailConfig - Updates Email Configuration
func UpdateEmailConfig(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Email - UpdateEmailConfig")
	defer util.CommonDeferred(w, r, "Controller", "Email", "UpdateEmailConfig")

	var reqMap data.EmailConfig

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.EmailID, reqMap.EmailName, reqMap.SMTPHost, reqMap.SMTPUser, reqMap.SMTPPassword, reqMap.Signature)
	ValidateFloat := ValidateNotNullStructFloat(reqMap.SMTPPort)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "EMAIL_CONFIG", nil, map[string]string{"email_id": reqMap.EmailID}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateEmailConfig(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetEmailConfigInfo - Gets Email Config Info
func GetEmailConfigInfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Email - GetEmailConfigInfo")
	defer util.CommonDeferred(w, r, "Controller", "Email", "GetEmailConfigInfo")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetEmailConfigInfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// EmailListing - Datatable Email listing with filter and order
func EmailListing(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Email - EmailListing")
	defer util.CommonDeferred(w, r, "Controller", "Email", "EmailListing")

	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.EmailListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
