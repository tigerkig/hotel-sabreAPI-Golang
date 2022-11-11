package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddMealType - Add Meal Type
func AddMealType(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Meal_Type - AddMealType")
	defer util.CommonDeferred(w, r, "Controller", "V_Meal_Type", "AddMealType")

	var reqMap data.MealType

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Type)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "MEAL_TYPE", map[string]string{"meal_type": reqMap.Type}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddMealType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateMealType - Update Meal Type
func UpdateMealType(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Meal_Type - UpdateMealType")
	defer util.CommonDeferred(w, r, "Controller", "V_Meal_Type", "UpdateMealType")

	var reqMap data.MealType

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Type)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "MEAL_TYPE", map[string]string{"meal_type": reqMap.Type}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateMealType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// MealTypeListing - Datatable Meal type listing with filter and order
func MealTypeListing(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Meal_Type - MealTypeListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Meal_Type", "MealTypeListing")

	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.MealTypeListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetMealTypeList - Get Meal Type List For Other Module
func GetMealTypeList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Meal_Type - GetMealTypeList")
	defer util.CommonDeferred(w, r, "Controller", "V_Meal_Type", "GetMealTypeList")

	retMap, err := model.GetMealTypeList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
