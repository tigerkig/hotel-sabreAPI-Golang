package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login - Admin Panel Login
func Login(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Controller - Login")
	defer util.CommonDeferred(w, r, "Controller", "Controller", "Login")
	var reqMap loginCredentials
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}
	var username = reqMap.Username
	var passwordTmp = reqMap.Password
	var password = util.GeneratePasswordHash(passwordTmp)
	retObj, flag := model.AdminLogin(r, username, password)
	if flag == 0 {
		util.Respond(r, w, nil, 401, "10001")
		return
	} else if flag == 2 {
		util.Respond(r, w, nil, 401, "10001")
		return
	}
	util.RespondData(r, w, retObj, 200)
}

// GetAuthDetails - It authenticate and returns user auth details
func GetAuthDetails(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "controller - controller - GetAuthDetails")
	defer util.CommonDeferred(w, r, "Controller", "controller", "GetAuthDetails")
	vars := mux.Vars(r)
	token := vars["token"]
	if token == "" {
		util.RespondBadRequest(r, w)
		return
	}
	SessionData, err := model.GetAuthDetails(r, token)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, SessionData, 200)
}

// Logout - Admin Panel Logout
func Logout(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Controller - Logout")
	defer util.CommonDeferred(w, r, "Controller", "Controller", "Logout")
	if model.Logout(r) {
		util.RespondData(r, w, nil, 200)
		return
	}
	util.Respond(r, w, nil, 200, "")
}

// PartnerLogin - Partner Panel Login
func PartnerLogin(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Controller - PartnerLogin")
	defer util.CommonDeferred(w, r, "Controller", "Controller", "PartnerLogin")
	var reqMap loginCredentials
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}
	var username = reqMap.Username
	var passwordTmp = reqMap.Password
	var password = util.GeneratePasswordHash(passwordTmp)
	retObj, flag := model.PartnerMultiHotelLogin(r, username, password)
	if flag == 0 {
		util.Respond(r, w, nil, 401, "10001")
		return
	} else if flag == 2 {
		util.Respond(r, w, nil, 401, "10001")
		return
	}
	util.RespondData(r, w, retObj, 200)
}

// GetConsoleAuthDetails - It authenticate and returns user auth details
func GetConsoleAuthDetails(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "controller - controller - GetConsoleAuthDetails")
	defer util.CommonDeferred(w, r, "Controller", "controller", "GetConsoleAuthDetails")
	vars := mux.Vars(r)
	token := vars["token"]
	if token == "" {
		util.RespondBadRequest(r, w)
		return
	}
	SessionData, err := model.GetConsoleAuthDetails(r, token)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}
	SessionDataInterface := make(map[string]interface{})
	if SessionData != nil {
		for k, v := range SessionData {
			SessionDataInterface[k] = v
		}
		SessionDataInterface["hotel_list"] = model.PartnerHotelList(r)
	}

	util.RespondData(r, w, SessionDataInterface, 200)
}
