package controller

import (
	"net/http"
	"tp-api-common/util"
	"tp-system/model"
)

// GetRoomTypeList - Get Room Type List For Other Module
func GetRoomTypeList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Controller - GetRoomTypeList")
	checkoutid := r.URL.Query().Get("hotelid")
	if checkoutid == "" {
		util.RespondBadRequest(r, w)
		return
	}
	Data, err := model.GetRoomTypeList(r, checkoutid)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, Data, 200)
}

// GetRatePlanFromRoom - Return rateplan including cancellation policy and other stuff by passing room id
func GetRatePlanFromRoom(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Controller - GetRatePlanFromRoom")
	hotelid := r.URL.Query().Get("hotelid")
	roomid := r.URL.Query().Get("roomid")
	if hotelid == "" || roomid == "" {
		util.RespondBadRequest(r, w)
		return
	}
	Data, err := model.GetRatePlanFromRoom(r, hotelid, roomid)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, Data, 200)
}

// GetHotelRoomRateData - Get All Data like Inv, Rate, Min, SS, CTA, CTD Of Room Type And Rate Plan Data Of Hotel
func GetHotelRoomRateData(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Controller - GetHotelRoomRateData")
	defer util.CommonDeferred(w, r, "Controller", "Controller", "GetHotelRoomRateData")
	var reqMap = make(map[string]interface{})
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	hotelid := r.URL.Query().Get("hotelid")
	roomid := r.URL.Query().Get("roomid")
	rateid := r.URL.Query().Get("rateid")

	if month == "" || year == "" || hotelid == "" || roomid == "" || rateid == "" {
		util.RespondBadRequest(r, w)
		return
	}

	reqMap["month"] = month
	reqMap["year"] = year
	reqMap["room_id"] = roomid
	reqMap["rate_id"] = rateid
	reqMap["hotel_id"] = hotelid

	Data, err := model.GetHotelRoomRateData(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")

}
