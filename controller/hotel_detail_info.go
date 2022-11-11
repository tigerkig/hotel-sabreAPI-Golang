package controller

import (
	"net/http"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// GetHotelDetailInfo - Get Hotel Detail Info
func GetHotelDetailInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Controller - GetHotelDetailInfo")
	vars := mux.Vars(r)
	checkoutid := vars["id"]
	if checkoutid == "" {
		util.RespondBadRequest(r, w)
		return
	}
	Data, err := model.GetHotelDetailInfo(r, checkoutid)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, Data, 200)
}
