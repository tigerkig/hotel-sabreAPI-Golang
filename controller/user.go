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

// AddUser - Add User
func AddUser(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User - AddUser")
	defer util.CommonDeferred(w, r, "Controller", "V_User", "AddUser")
	var reqMap data.User

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	var username = reqMap.Username
	var password = reqMap.Password
	var name = reqMap.Name
	var email = reqMap.Email
	var cc = reqMap.CC
	var phoneCode = reqMap.PhoneCode
	var phone = reqMap.Phone
	var role = reqMap.Role
	var privileges = reqMap.Privileges

	var NewUserRole = reqMap.NewUserRole

	ValidateString := ValidateNotNullStructString(username, password, name, privileges, email, cc, role, reqMap.BirthDate, reqMap.Address)
	ValidateFloat := ValidateNotNullStructFloat(phoneCode, phone)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// 2020-06-20 - HK - On Fly User Role Adding Validation Added
	if NewUserRole != "" {
		cnt, err := model.CheckDuplicateRecords(r, "ROLE", nil, map[string]string{"role": NewUserRole}, "0")
		if err != nil || cnt != 0 {
			if err != nil {
				util.RespondBadRequest(r, w)
			} else {
				util.Respond(r, w, nil, 409, "10010")
			}
			return
		}
	}
	// 2020-06-20 - HK - On Fly User Role Adding Validation Added

	cnt, err := model.CheckDuplicateRecords(r, "USER", nil, map[string]string{"username": username}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10006")
		}
		return
	}

	cnt, err = model.CheckDuplicateRecords(r, "USER_PROFILE", map[string]string{"phone_code": fmt.Sprintf("%.0f", phoneCode), "phone": fmt.Sprintf("%.0f", phone)}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10008")
		}
		return
	}

	cnt, err = model.CheckDuplicateRecords(r, "USER_PROFILE", nil, map[string]string{"email": email}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10007")
		}
		return
	}

	// 2020-06-20 - HK - On Fly User Role Added
	if NewUserRole != "" {
		retMap, err := model.OnflyAddUserRole(r, NewUserRole, privileges)
		if err != nil {
			util.RespondWithError(r, w, "500")
			return
		}
		reqMap.Role = retMap["id"].(string)
	}
	// 2020-06-20 - HK - On Fly User Role Added

	flag := model.AddUser(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateUser - Update User
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User - UpdateUser")
	defer util.CommonDeferred(w, r, "Controller", "V_User", "UpdateUser")
	var reqMap data.User

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.ID = id
	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Username, reqMap.Name, reqMap.Email, reqMap.CC, reqMap.BirthDate, reqMap.Address) // 2020-06-20 - HK reqMap.Role, reqMap.Privileges Moved To Reset Screen
	ValidateFloat := ValidateNotNullStructFloat(reqMap.PhoneCode, reqMap.Phone)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "USER_PROFILE", map[string]string{"phone_code": fmt.Sprintf("%.0f", reqMap.PhoneCode), "phone": fmt.Sprintf("%.0f", reqMap.Phone)}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10008")
		}
		return
	}

	cnt, err = model.CheckDuplicateRecords(r, "USER_PROFILE", nil, map[string]string{"email": reqMap.Email}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10007")
		}
		return
	}

	flag := model.UpdateUser(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetUser - Get User Detail By ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User - GetUser")
	defer util.CommonDeferred(w, r, "Controller", "V_User", "GetUser")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetUserInfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// UserListing - User Listing
func UserListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User - UserListing")
	defer util.CommonDeferred(w, r, "Controller", "V_User", "UserListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.UserListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// ResetPwdAndPrivileges - Reset Password And Privileges Of User
func ResetPwdAndPrivileges(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User - ResetPwdAndPrivileges")
	defer util.CommonDeferred(w, r, "Controller", "V_User", "ResetPwdAndPrivileges")
	var reqMap data.ResetPwdAndPrivileges

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.ID = id

	var role = reqMap.Role
	var privileges = reqMap.Privileges

	ValidateString := ValidateNotNullStructString(reqMap.ID, role, privileges)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// 2020-06-22 - HK - On Fly User Role Adding Validation Added
	var NewUserRole = reqMap.NewUserRole
	if NewUserRole != "" {
		cnt, err := model.CheckDuplicateRecords(r, "ROLE", nil, map[string]string{"role": NewUserRole}, "0")
		if err != nil || cnt != 0 {
			if err != nil {
				util.RespondBadRequest(r, w)
			} else {
				util.Respond(r, w, nil, 409, "10010")
			}
			return
		}
		retMap, err := model.OnflyAddUserRole(r, NewUserRole, privileges)
		if err != nil {
			util.RespondWithError(r, w, "500")
			return
		}
		reqMap.Role = retMap["id"].(string)
	}
	// 2020-06-22 - HK - On Fly User Role Adding Validation Added

	flag := model.ResetPwdAndPrivileges(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetUserList - Get User Active List For Other Module Use
func GetUserList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User - GetUserList")
	defer util.CommonDeferred(w, r, "Controller", "V_User", "GetUserList")
	retMap, err := model.GetUserList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// ResetPassword - Reset Password Of Admin And Partner
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_User - ResetPassword")
	defer util.CommonDeferred(w, r, "Controller", "V_User", "ResetPassword")
	var reqMap data.ResetPassword
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.OldPassword, reqMap.NewPassword)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	if !model.CheckOldPassword(r, reqMap) {
		util.Respond(r, w, nil, 406, "Invalid old password!!")
		return
	}

	flag := model.ResetPassword(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}
