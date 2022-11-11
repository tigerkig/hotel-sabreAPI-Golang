package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddImageCategory - Add Image Category
func AddImageCategory(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Image_Category - AddImageCategory")
	defer util.CommonDeferred(w, r, "Controller", "V_Image_Category", "AddImageCategory")
	var reqMap data.ImageCategory

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Category)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "IMAGE_CATEGORY", nil, map[string]string{"name": reqMap.Category}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddImageCategory(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateImageCategory - Update Image Category
func UpdateImageCategory(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Image_Category - UpdateImageCategory")
	defer util.CommonDeferred(w, r, "Controller", "V_Image_Category", "UpdateImageCategory")
	var reqMap data.ImageCategory

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Category)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "IMAGE_CATEGORY", nil, map[string]string{"name": reqMap.Category}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateImageCategory(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// ImageCategoryListing - Datatable Image Category listing with filter and order
func ImageCategoryListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Image_Category - ImageCategoryListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Image_Category", "ImageCategoryListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.ImageCategoryListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetImageCategoryList - Get Image Category List For Other Module
func GetImageCategoryList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Image_Category - GetImageCategoryList")
	defer util.CommonDeferred(w, r, "Controller", "V_Image_Category", "GetImageCategoryList")
	retMap, err := model.GetImageCategoryList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
