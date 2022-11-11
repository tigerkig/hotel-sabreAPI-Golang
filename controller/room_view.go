package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddRoomView - Add Room View
func AddRoomView(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Room_View - AddRoomView")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_View", "AddRoomView")
	var reqMap data.RoomView

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.RoomView)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "ROOM_VIEW", nil, map[string]string{"room_view_name": reqMap.RoomView}, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddRoomView(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateRoomView - Update Room View
func UpdateRoomView(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Room_View - UpdateRoomView")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_View", "UpdateRoomView")
	var reqMap data.RoomView

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.RoomView)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "ROOM_VIEW", nil, map[string]string{"room_view_name": reqMap.RoomView}, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateRoomView(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetRoomViewList - Get Room View List For Other Module
func GetRoomViewList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Room_View - GetRoomViewList")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_View", "GetRoomViewList")
	retMap, err := model.GetRoomViewList(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// RoomViewListing - Datatable Room View listing with filter and order
func RoomViewListing(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Room_View - RoomViewListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_View", "RoomViewListing")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.RoomViewListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetRoomView - Get Room View Detail By ID - 2021-04-21 - HK
func GetRoomView(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Room_View - GetRoomView")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_View", "GetRoomView")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := model.GetRoomView(r, id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}
