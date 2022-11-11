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

// AddPopularCity - Add Popular City
func AddPopularCity(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Popular_City - AddPopularCity")
	defer util.CommonDeferred(w, r, "Controller", "Popular_City", "AddPopularCity")

	var reqMap = make(map[string]interface{})

	CityID := r.FormValue("city_id")
	SortOrder := r.FormValue("sort_order")
	CityName := r.FormValue("city_name")
	Description := r.FormValue("description")

	if CityID == "" || SortOrder == "" || CityName == "" || Description == "" {
		util.RespondBadRequest(r, w)
		return
	}

	imgfile, _, err := r.FormFile("image")
	if err != nil || imgfile == nil {
		util.RespondBadRequest(r, w)
		return
	}

	/*
		    // 2021-04-26 - HK - START
			// Purpose : Use of this code was to stop adding more than 9 records as per our design template.
			// Now it's commented because as per user flexibility, user should not be restricted to add
			// more than 9 records. But to keep our design template as it is designed, it should be handled on
			// status change event. hence this checking has been moved to status change event.
			TotalRecordCount, err := model.GetPopularCityCount(r)
			if err != nil {
				util.RespondBadRequest(r, w)
				return
			}

			if TotalRecordCount >= 9 {
				util.Respond(r, w, nil, 406, "100014")
				return
			}
			// 2021-04-26 - HK - END
	*/

	cnt, err := model.CheckDuplicateRecords(r, "POPULARCITY", map[string]string{"sort_order": SortOrder, "city_id": CityID}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	reqMap["city_id"] = CityID
	reqMap["city_name"] = CityName
	reqMap["sort_order"] = SortOrder
	reqMap["description"] = Description

	if imgfile != nil {
		ImageName, err := UploadSingleImageFormData(r, "popular_city", "image")
		if err != nil {
			util.RespondBadRequest(r, w)
			return
		}
		reqMap["image"] = ImageName
		defer imgfile.Close()
	}

	flag := model.AddPopularCity(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdatePopularCity - Add Popular City
func UpdatePopularCity(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Popular_City - UpdatePopularCity")
	defer util.CommonDeferred(w, r, "Controller", "Popular_City", "UpdatePopularCity")

	var reqMap = make(map[string]interface{})

	vars := mux.Vars(r)
	ID := vars["id"]
	CityID := r.FormValue("city_id")
	SortOrder := r.FormValue("sort_order")
	CityName := r.FormValue("city_name")
	Description := r.FormValue("description")

	if ID == "" || CityID == "" || SortOrder == "" || CityName == "" || Description == "" {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "POPULARCITY", map[string]string{"sort_order": SortOrder, "city_id": CityID}, nil, ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	reqMap["city_id"] = CityID
	reqMap["city_name"] = CityName
	reqMap["sort_order"] = SortOrder
	reqMap["description"] = Description
	reqMap["id"] = ID

	imgfile, _, err := r.FormFile("image")
	if imgfile != nil {
		ImageName, err := UploadSingleImageFormData(r, "popular_city", "image")
		if err != nil {
			util.LogIt(r, fmt.Sprint("Controller - Popular_City - UpdatePopularCity - Error While Updating Popular City Pic"))
			util.RespondBadRequest(r, w)
			return
		}
		reqMap["image"] = ImageName
	} else {
		reqMap["image"] = r.FormValue("image")
	}

	flag := model.UpdatePopularCity(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// PopularCityListing - Return Datatable Listing Of Popular City
func PopularCityListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Popular_City - PopularCityListing")
	defer util.CommonDeferred(w, r, "Controller", "Popular_City", "PopularCityListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.PopularCityListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetPopularCityinfo - Get Popular City Detail By ID
func GetPopularCityinfo(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - Popular_City - GetPopularCityinfo")
	defer util.CommonDeferred(w, r, "Controller", "Popular_City", "GetPopularCityinfo")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetPopularCityinfo(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// UpdatePopularCityStatus - Updates Status Of Popular City - 2021-04-26 - HK
func UpdatePopularCityStatus(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Popular_City - UpdatePopularCityStatus")
	defer util.CommonDeferred(w, r, "Controller", "Popular_City", "UpdatePopularCityStatus")

	var reqMap data.StatusForSingle
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.ID = id
	var status = int(reqMap.Status)

	validateFlag := ValidateNotNullStructString(id)
	if validateFlag == 0 || id == "" {
		util.RespondBadRequest(r, w)
		return
	}

	if status == 1 {
		TotalRecordCount, err := model.GetPopularCityCount(r, true)
		if err != nil {
			util.Respond(r, w, nil, 500, "100014")
			return
		}

		if TotalRecordCount >= 9 {
			util.Respond(r, w, nil, 406, "100014")
			return
		}
	}

	_, code, err1 := model.UpdatePopularCityStatus(r, "POPULARCITY", status, id)
	if err1 != nil {
		util.Respond(r, w, nil, code, "100014")
		return
	}

	util.Respond(r, w, nil, code, "")
}
