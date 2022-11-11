package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

func UpdateHotelStatus(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Common - V_Common")
	defer util.CommonDeferred(w, r, "Controller", "V_Common", "V_Common")
	var reqMap data.StatusForSingle
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	propertyStatus := vars["propertyStatus"]
	prStatus, _ := strconv.Atoi(propertyStatus)
	reqMap.ID = id
	var status = int(reqMap.Status)

	if reqMap.Status != 7 && reqMap.Status != 2 && reqMap.Status != 1 {
		util.RespondBadRequest(r, w)
		return
	}

	flag, httpStatus, err := model.UpdateStatusModuleWise(r, "HOTEL", status, reqMap.ID)
	if flag == 0 || util.CheckErrorLog(r, err) {
		util.Respond(r, w, nil, httpStatus, "100014")
		return
	}

	switch status {
	// Status active
	case 1:
		if prStatus == 1 {
			model.UpdateSearchList(reqMap.ID, "List")
			util.SysLogIt("UpdateSearchList For Hotel start")
			model.CacheChn <- model.CacheObj{
				Type: "updateHotelWithProperty",
				ID:   id,
			}
			util.SysLogIt("UpdateSearchList For Hotel end")

		} else {
			model.UpdateHotelStatusToLive(r, prStatus, reqMap.ID)
		}
		model.UpdateHotelierLoginStatus(r, reqMap.ID, 1)
	//Status inactive/blacklisted
	case 2:
		model.UpdateHotelierLoginStatus(r, reqMap.ID, 2)
	case 7:
		model.UpdateHotelStatusToLive(r, 2, reqMap.ID)
	}

	util.Respond(r, w, nil, 204, "")
}
