package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddCms - Add CMS
func AddCms(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Cms - AddCms")
	defer util.CommonDeferred(w, r, "Controller", "Cms", "AddCms")
	var reqMap data.Cms

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	// Get Slug for Newly Added Page
	reqMap.Slug = model.Slug(reqMap.Page)

	ValidateString := ValidateNotNullStructString(reqMap.Page, reqMap.SortCode, reqMap.Slug, reqMap.Content)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.AddCms(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateCms - Update Cms
func UpdateCms(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Cms - UpdateCms")
	defer util.CommonDeferred(w, r, "Controller", "Cms", "UpdateCms")
	var reqMap data.Cms

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Page, reqMap.Content)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := model.UpdateCms(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetCms -  Return Cms Details
func GetCms(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Cms - GetCms")
	defer util.CommonDeferred(w, r, "Controller", "Cms", "GetCms")
	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetCms(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// CmsListing - Return Datatable Listing Of Cms
func CmsListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Cms - CmsListing")
	defer util.CommonDeferred(w, r, "Controller", "Cms", "CmsListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.CmsListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
