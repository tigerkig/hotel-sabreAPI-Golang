package controller

import (
	"encoding/json"
	"net/http"
	"tp-system/model"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/mux"
)

// GetPrivilegeList - Get List Of Privilege Module Wise
func GetPrivilegeList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "controller - controller - GetPrivilegeList")
	defer util.CommonDeferred(w, r, "Controller", "controller", "GetPrivilegeList")
	stuff := make(map[string]interface{})
	List, err := model.GetPrivilegeList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	if len(List) == 0 {
		stuff["data"] = []string{}
	} else {
		stuff["data"] = List
	}

	util.RespondData(r, w, stuff, 200)
}

// AddUserRole - Add User Role
func AddUserRole(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User_Role - AddUserRole")
	defer util.CommonDeferred(w, r, "Controller", "V_User_Role", "AddUserRole")
	var reqMap data.UserRole

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	var role = reqMap.Role
	var privileges = reqMap.Privileges

	ValidateString := ValidateNotNullStructString(role, privileges)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "ROLE", nil, map[string]string{"role": reqMap.Role}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddUserRole(r, role, privileges)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// GetUserRole - Get User Role Detail By ID
func GetUserRole(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User_Role - GetUserRole")
	defer util.CommonDeferred(w, r, "Controller", "V_User_Role", "GetUserRole")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetUserRole(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// UpdateUserRole - Update User Role
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User_Role - UpdateUserRole")
	defer util.CommonDeferred(w, r, "Controller", "V_User_Role", "UpdateUserRole")
	var reqMap data.UserRole

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var role = reqMap.Role
	var privileges = reqMap.Privileges

	ValidateString := ValidateNotNullStructString(id, role, privileges)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "ROLE", nil, map[string]string{"role": reqMap.Role}, id)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateUserRole(r, id, role, privileges)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetUserRoleList - Get User Role Active List For Other Module Use
func GetUserRoleList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User_Role - GetUserRoleList")
	defer util.CommonDeferred(w, r, "Controller", "V_User_Role", "GetUserRoleList")
	retMap, err := model.GetUserRoleList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// UserRoleListing - User Role Listing
func UserRoleListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User_Role - UserRoleListing")
	defer util.CommonDeferred(w, r, "Controller", "V_User_Role", "UserRoleListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.UserRoleListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
